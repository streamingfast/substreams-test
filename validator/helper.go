package validator

import "github.com/streamingfast/substreams-test/validator/config"

func (v *Validator) shouldIgnoreEntity(entityName string) bool {
	if entity, found := v.config[entityName]; found {
		return entity.Ignore
	}
	return false
}

func (v *Validator) getGraphQLFieldName(entityName string, fieldName string) string {
	if entity, found := v.config[entityName]; found {
		if field, ok := entity.Fields[fieldName]; ok {
			if field.Rename != "" {
				fieldName = field.Rename
			}
		}
	}
	return fieldName
}

func (v *Validator) shouldIgnoreField(entityName string, fieldName string) bool {
	if entity, found := v.config[entityName]; found {
		if field, ok := entity.Fields[fieldName]; ok {
			return field.Ignore
		}
	}
	return false
}

func (v *Validator) isGraphQLAssociatedField(entityName string, fieldName string) bool {
	if entity, found := v.config[entityName]; found {
		if field, ok := entity.Fields[fieldName]; ok {
			return field.Association
		}
	}
	return false
}

func (v *Validator) isGraphQLArrayField(entityName string, fieldName string) bool {
	if entity, found := v.config[entityName]; found {
		if field, ok := entity.Fields[fieldName]; ok {
			return field.Array
		}
	}
	return false
}

func (v *Validator) getFieldOpt(entityName string, fieldName string) config.Options {
	if entity, found := v.config[entityName]; found {
		if field, ok := entity.Fields[fieldName]; ok {
			if v.defaultOptions.Error != nil && field.Opt.Error == nil {
				field.Opt.Error = v.defaultOptions.Error
			}
			return field.Opt
		}
	}
	return v.defaultOptions
}
