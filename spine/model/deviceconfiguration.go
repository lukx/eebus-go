package model

type DeviceConfigurationKeyIdType uint

type DeviceConfigurationKeyValueStringType string

type DeviceConfigurationKeyNameType string

const (
	DeviceConfigurationKeyNameTypePeakPowerOfPVSystem         DeviceConfigurationKeyNameType = "peakPowerOfPvSystem"
	DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor    DeviceConfigurationKeyNameType = "pvCurtailmentLimitFactor"
	DeviceConfigurationKeyNameTypeAsymmetricChargingSupported DeviceConfigurationKeyNameType = "asymmetricChargingSupported"
	DeviceConfigurationKeyNameTypeCommunicationsStandard      DeviceConfigurationKeyNameType = "communicationsStandard"
)

type DeviceConfigurationKeyValueTypeType string

const (
	DeviceConfigurationKeyValueTypeTypeBoolean      DeviceConfigurationKeyValueTypeType = "boolean"
	DeviceConfigurationKeyValueTypeTypeDate         DeviceConfigurationKeyValueTypeType = "date"
	DeviceConfigurationKeyValueTypeTypeDateTime     DeviceConfigurationKeyValueTypeType = "dateTime"
	DeviceConfigurationKeyValueTypeTypeDuration     DeviceConfigurationKeyValueTypeType = "duration"
	DeviceConfigurationKeyValueTypeTypeString       DeviceConfigurationKeyValueTypeType = "string"
	DeviceConfigurationKeyValueTypeTypeTime         DeviceConfigurationKeyValueTypeType = "time"
	DeviceConfigurationKeyValueTypeTypeScaledNumber DeviceConfigurationKeyValueTypeType = "scaledNumber"
)

type DeviceConfigurationKeyValueValueType struct {
	Boolean      *bool                                  `json:"boolean,omitempty"`
	Date         *DateType                              `json:"date,omitempty"`
	DateTime     *DateTimeType                          `json:"dateTime,omitempty"`
	Duration     *DurationType                          `json:"duration,omitempty"`
	String       *DeviceConfigurationKeyValueStringType `json:"string,omitempty"`
	Time         *TimeType                              `json:"time,omitempty"`
	ScaledNumber *ScaledNumberType                      `json:"scaledNumber,omitempty"`
}

type DeviceConfigurationKeyValueValueElementsType struct {
	Boolean      *ElementTagType           `json:"boolean,omitempty"`
	Date         *ElementTagType           `json:"date,omitempty"`
	DateTime     *ElementTagType           `json:"dateTime,omitempty"`
	Duration     *ElementTagType           `json:"duration,omitempty"`
	String       *ElementTagType           `json:"string,omitempty"`
	Time         *ElementTagType           `json:"time,omitempty"`
	ScaledNumber *ScaledNumberElementsType `json:"scaledNumber,omitempty"`
}

type DeviceConfigurationKeyValueDataType struct {
	KeyId             *DeviceConfigurationKeyIdType         `json:"keyId,omitempty" eebus:"key"`
	Value             *DeviceConfigurationKeyValueValueType `json:"value,omitempty"`
	IsValueChangeable *bool                                 `json:"isValueChangeable,omitempty"`
}

type DeviceConfigurationKeyValueDataElementsType struct {
	KeyId             *ElementTagType                               `json:"keyId,omitempty"`
	Value             *DeviceConfigurationKeyValueValueElementsType `json:"value,omitempty"`
	IsValueChangeable *ElementTagType                               `json:"isValueChangeable,omitempty"`
}

type DeviceConfigurationKeyValueListDataType struct {
	DeviceConfigurationKeyValueData []DeviceConfigurationKeyValueDataType `json:"deviceConfigurationKeyValueData,omitempty"`
}

type DeviceConfigurationKeyValueListDataSelectorsType struct {
	KeyId *DeviceConfigurationKeyIdType `json:"keyId,omitempty"`
}

type DeviceConfigurationKeyValueDescriptionDataType struct {
	KeyId       *DeviceConfigurationKeyIdType        `json:"keyId,omitempty" eebus:"key"`
	KeyName     *DeviceConfigurationKeyNameType      `json:"keyName,omitempty"`
	ValueType   *DeviceConfigurationKeyValueTypeType `json:"valueType,omitempty"`
	Unit        *UnitOfMeasurementType               `json:"unit,omitempty"`
	Label       *LabelType                           `json:"label,omitempty"`
	Description *DescriptionType                     `json:"description,omitempty"`
}

type DeviceConfigurationKeyValueDescriptionDataElementsType struct {
	KeyId       *ElementTagType `json:"keyId,omitempty"`
	KeyName     *ElementTagType `json:"keyName,omitempty"`
	ValueType   *ElementTagType `json:"valueType,omitempty"`
	Unit        *ElementTagType `json:"unit,omitempty"`
	Label       *ElementTagType `json:"label,omitempty"`
	Description *ElementTagType `json:"description,omitempty"`
}

type DeviceConfigurationKeyValueDescriptionListDataType struct {
	DeviceConfigurationKeyValueDescriptionData []DeviceConfigurationKeyValueDescriptionDataType `json:"deviceConfigurationKeyValueDescriptionData,omitempty"`
}

type DeviceConfigurationKeyValueDescriptionListDataSelectorsType struct {
	KeyId   *DeviceConfigurationKeyIdType `json:"keyId,omitempty"`
	KeyName *string                       `json:"keyName,omitempty"`
}

type DeviceConfigurationKeyValueConstraintsDataType struct {
	KeyId         *DeviceConfigurationKeyIdType         `json:"keyId,omitempty" eebus:"key"`
	ValueRangeMin *DeviceConfigurationKeyValueValueType `json:"valueRangeMin,omitempty"`
	ValueRangeMax *DeviceConfigurationKeyValueValueType `json:"valueRangeMax,omitempty"`
	ValueStepSize *DeviceConfigurationKeyValueValueType `json:"valueStepSize,omitempty"`
}

type DeviceConfigurationKeyValueConstraintsDataElementsType struct {
	KeyId         *ElementTagType                               `json:"keyId,omitempty"`
	ValueRangeMin *DeviceConfigurationKeyValueValueElementsType `json:"valueRangeMin,omitempty"`
	ValueRangeMax *DeviceConfigurationKeyValueValueElementsType `json:"valueRangeMax,omitempty"`
	ValueStepSize *DeviceConfigurationKeyValueValueElementsType `json:"valueStepSize,omitempty"`
}

type DeviceConfigurationKeyValueConstraintsListDataType struct {
	DeviceConfigurationKeyValueConstraintsData []DeviceConfigurationKeyValueConstraintsDataType `json:"deviceConfigurationKeyValueConstraintsData,omitempty"`
}

type DeviceConfigurationKeyValueConstraintsListDataSelectorsType struct {
	KeyId *DeviceConfigurationKeyIdType `json:"keyId,omitempty"`
}
