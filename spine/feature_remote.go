package spine

import (
	"fmt"
	"time"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/rickb777/date/period"
)

const defaultMaxResponseDelay = time.Duration(time.Second * 10)

type FeatureRemoteImpl struct {
	*FeatureImpl
	entity           *EntityRemoteImpl
	functionDataMap  map[model.FunctionType]FunctionData
	maxResponseDelay *time.Duration
}

func NewFeatureRemoteImpl(id uint, entity *EntityRemoteImpl, ftype model.FeatureTypeType, role model.RoleType) *FeatureRemoteImpl {
	res := &FeatureRemoteImpl{
		FeatureImpl: NewFeatureImpl(
			featureAddressType(id, entity.Address()),
			ftype,
			role),
		entity:          entity,
		functionDataMap: make(map[model.FunctionType]FunctionData),
	}
	for _, fd := range CreateFunctionData[FunctionData](ftype) {
		res.functionDataMap[fd.Function()] = fd
	}

	res.operations = make(map[model.FunctionType]*Operations)

	return res
}

func (r *FeatureRemoteImpl) Data(function model.FunctionType) any {
	return r.functionData(function).DataAny()
}

func (r *FeatureRemoteImpl) UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType) {
	r.functionData(function).UpdateDataAny(data, filterPartial, filterDelete)
	// TODO: fire event
}

func (r *FeatureRemoteImpl) Sender() Sender {
	return r.Device().Sender()
}

func (r *FeatureRemoteImpl) Device() *DeviceRemoteImpl {
	return r.entity.Device()
}

func (r *FeatureRemoteImpl) Entity() *EntityRemoteImpl {
	return r.entity
}

func (r *FeatureRemoteImpl) SetOperations(functions []model.FunctionPropertyType) {
	r.operations = make(map[model.FunctionType]*Operations)
	for _, sf := range functions {
		r.operations[*sf.Function] = NewOperations(sf.PossibleOperations.Read != nil, sf.PossibleOperations.Write != nil)
	}
}

func (r *FeatureRemoteImpl) SetMaxResponseDelay(delay *model.MaxResponseDelayType) {
	if delay == nil {
		return
	}
	p, err := period.Parse(string(*delay))
	if err != nil {
		r.maxResponseDelay = util.Ptr(p.DurationApprox())
	} else {
		logging.Log.Debug(err)
	}
}

func (r *FeatureRemoteImpl) MaxResponseDelayDuration() time.Duration {
	if r.maxResponseDelay != nil {
		return *r.maxResponseDelay
	}
	return defaultMaxResponseDelay
}

func (r *FeatureRemoteImpl) functionData(function model.FunctionType) FunctionData {
	fd, found := r.functionDataMap[function]
	if !found {
		panic(fmt.Errorf("Data was not found for function '%s'", function))
	}
	return fd
}
