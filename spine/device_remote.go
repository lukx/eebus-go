package spine

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine/model"
)

type DeviceRemoteImpl struct {
	*DeviceImpl

	ski string

	entities      []*EntityRemoteImpl
	entitiesMutex sync.Mutex

	sender Sender

	localDevice *DeviceLocalImpl

	// Heartbeat Sender
	heartbeatSender *HeartbeatSender
}

var _ SpineDataProcessing = (*DeviceRemoteImpl)(nil)

func NewDeviceRemoteImpl(localDevice *DeviceLocalImpl, ski string, writeHandler SpineDataConnection) *DeviceRemoteImpl {
	sender := NewSender(writeHandler)
	res := DeviceRemoteImpl{
		DeviceImpl:      NewDeviceImpl(nil, nil, nil),
		ski:             ski,
		localDevice:     localDevice,
		sender:          sender,
		heartbeatSender: NewHeartbeatSender(sender),
	}
	res.addNodeManagement()

	return &res
}

// return the device SKI
func (d *DeviceRemoteImpl) Ski() string {
	return d.ski
}

// Needs to be called by the CEM implementation once a subscription for the local DeviceDiagnosis server feature is received
func (d *DeviceRemoteImpl) StartHeartbeatSend(senderAddr, destinationAddr *model.FeatureAddressType) {
	d.heartbeatSender.StartHeartbeatSend(senderAddr, destinationAddr)
}

func (d *DeviceRemoteImpl) IsHeartbeatMsgCounter(msgCounter model.MsgCounterType) bool {
	return d.heartbeatSender.IsHeartbeatMsgCounter(msgCounter)
}

// Needs to be called by the CEM implementation once a subscription for the local DeviceDiagnosis server feature is removed
func (d *DeviceRemoteImpl) Stopheartbeat() {
	d.heartbeatSender.StopHeartbeat()
}

// this connection is closed
func (d *DeviceRemoteImpl) CloseConnection() {
	d.heartbeatSender.StopHeartbeat()
}

// processing incoming SPINE message from the associated SHIP connection
func (d *DeviceRemoteImpl) HandleIncomingSpineMesssage(message []byte) (*model.MsgCounterType, error) {
	datagram := model.Datagram{}
	if err := json.Unmarshal([]byte(message), &datagram); err != nil {
		return nil, err
	}
	err := d.localDevice.ProcessCmd(datagram.Datagram, d)
	if err != nil {
		logging.Log.Trace(err)
	}

	return datagram.Datagram.Header.MsgCounter, nil
}

func (d *DeviceRemoteImpl) addNodeManagement() {
	deviceInformation := d.addNewEntity(model.EntityTypeTypeDeviceInformation, NewAddressEntityType([]uint{DeviceInformationEntityId}))
	nodeManagement := NewFeatureRemoteImpl(deviceInformation.NextFeatureId(), deviceInformation, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	deviceInformation.AddFeature(nodeManagement)
}

func (d *DeviceRemoteImpl) Sender() Sender {
	return d.sender
}

// Return an entity with a given address
func (d *DeviceRemoteImpl) Entity(id []model.AddressEntityType) *EntityRemoteImpl {
	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	for _, e := range d.entities {
		if reflect.DeepEqual(id, e.Address().Entity) {
			return e
		}
	}
	return nil
}

// Return all entities of this device
func (d *DeviceRemoteImpl) Entities() []*EntityRemoteImpl {
	return d.entities
}

// Return the feature for a given address
func (d *DeviceRemoteImpl) FeatureByAddress(address *model.FeatureAddressType) *FeatureRemoteImpl {
	entity := d.Entity(address.Entity)
	if entity != nil {
		return entity.Feature(address.Feature)
	}
	return nil
}

// Remove an entity with a given address from this device
func (d *DeviceRemoteImpl) RemoveByAddress(addr []model.AddressEntityType) *EntityRemoteImpl {
	entityForRemoval := d.Entity(addr)
	if entityForRemoval == nil {
		return nil
	}

	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	var newEntities []*EntityRemoteImpl
	for _, item := range d.entities {
		if !reflect.DeepEqual(item, entityForRemoval) {
			newEntities = append(newEntities, item)
		}
	}
	d.entities = newEntities

	return entityForRemoval
}

// Get the feature for a given entity, feature type and feature role
func (r *DeviceRemoteImpl) FeatureByEntityTypeAndRole(entity *EntityRemoteImpl, featureType model.FeatureTypeType, role model.RoleType) *FeatureRemoteImpl {
	if len(r.entities) < 1 {
		return nil
	}

	r.entitiesMutex.Lock()
	defer r.entitiesMutex.Unlock()

	for _, e := range r.entities {
		if entity != e {
			continue
		}
		for _, feature := range entity.Features() {
			if feature.Type() == featureType && feature.Role() == role {
				return feature
			}
		}
	}

	return nil
}

func (d *DeviceRemoteImpl) UpdateDevice(description *model.NetworkManagementDeviceDescriptionDataType) {
	if description != nil {
		if description.DeviceAddress != nil && description.DeviceAddress.Device != nil {
			d.address = description.DeviceAddress.Device
		}
		if description.DeviceType != nil {
			d.dType = description.DeviceType
		}
		if description.NetworkFeatureSet != nil {
			d.featureSet = description.NetworkFeatureSet
		}
	}
}

func (d *DeviceRemoteImpl) AddEntityAndFeatures(initialData bool, data *model.NodeManagementDetailedDiscoveryDataType) ([]*EntityRemoteImpl, error) {
	rEntites := make([]*EntityRemoteImpl, 0)

	for _, ei := range data.EntityInformation {
		if err := d.CheckEntityInformation(initialData, ei); err != nil {
			return nil, err
		}

		entityAddress := ei.Description.EntityAddress.Entity

		entity := d.Entity(entityAddress)
		if entity == nil {
			entity = d.addNewEntity(*ei.Description.EntityType, entityAddress)
			rEntites = append(rEntites, entity)
		}

		entity.SetDescription(ei.Description.Description)
		entity.RemoveAllFeatures()

		for _, fi := range data.FeatureInformation {
			if reflect.DeepEqual(fi.Description.FeatureAddress.Entity, entityAddress) {
				if f := unmarshalFeature(entity, fi); f != nil {
					entity.AddFeature(f)
				}
			}
		}

		// TOV-TODO: check this approach
		// if err := f.announceFeatureDiscovery(entity); err != nil {
		// 	return err
		// }
	}

	return rEntites, nil
}

// check if the provided entity information is correct
// provide initialData to check if the entity is new and not an update
func (d *DeviceRemoteImpl) CheckEntityInformation(initialData bool, entity model.NodeManagementDetailedDiscoveryEntityInformationType) error {
	description := entity.Description
	if description == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description")
	}

	if description.EntityAddress == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress")
	}

	if description.EntityAddress.Entity == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress.Entity")
	}

	// Consider on initial NodeManagement Detailed Discovery, the device being empty as it is not yet known
	if initialData {
		return nil
	}

	address := d.Address()
	if description.EntityAddress.Device != nil && address != nil && *description.EntityAddress.Device != *address {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: device address mismatch")
	}

	return nil
}

func (d *DeviceRemoteImpl) addNewEntity(eType model.EntityTypeType, address []model.AddressEntityType) *EntityRemoteImpl {
	newEntity := NewEntityRemoteImpl(d, eType, address)
	return d.addEntity(newEntity)
}

func (d *DeviceRemoteImpl) addEntity(entity *EntityRemoteImpl) *EntityRemoteImpl {
	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	d.entities = append(d.entities, entity)

	return entity
}

func unmarshalFeature(entity *EntityRemoteImpl,
	featureData model.NodeManagementDetailedDiscoveryFeatureInformationType,
) *FeatureRemoteImpl {
	var result *FeatureRemoteImpl

	if fid := featureData.Description; fid != nil {

		result = NewFeatureRemoteImpl(uint(*fid.FeatureAddress.Feature), entity, *fid.FeatureType, *fid.Role)

		result.SetDescription(fid.Description)
		result.SetMaxResponseDelay(fid.MaxResponseDelay)
		result.SetOperations(fid.SupportedFunction)
	}

	return result
}
