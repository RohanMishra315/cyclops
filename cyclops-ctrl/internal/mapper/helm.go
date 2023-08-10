package mapper

import (
	"sort"

	"github.com/cyclops-ui/cycops-ctrl/internal/models"
	"github.com/cyclops-ui/cycops-ctrl/internal/models/helm"
)

func HelmSchemaToFields(schema helm.Property) []models.Field {
	fields := make([]models.Field, 0, len(schema.Properties))

	for name, property := range schema.Properties {
		if property.Type == "array" {
			fields = append(fields, models.Field{
				Name:        name,
				Description: property.Description,
				Type:        mapHelmPropertyTypeToFieldType(property),
				DisplayName: mapTitle(name, property),
				ManifestKey: name,
				Items:       arrayItem(property.Items),
				Enum:        property.Enum,
			})
			continue
		}

		fields = append(fields, models.Field{
			Name:        name,
			Description: property.Description,
			Type:        mapHelmPropertyTypeToFieldType(property),
			DisplayName: mapTitle(name, property),
			ManifestKey: name,
			Properties:  HelmSchemaToFields(property),
			Enum:        property.Enum,
		})
	}

	return sortFields(fields, schema.Order)
}

func sortFields(fields []models.Field, order []string) []models.Field {
	ordersMap := make(map[string]int)

	for i, s := range order {
		ordersMap[s] = i
	}

	sort.Slice(fields, func(i, j int) bool {
		return ordersMap[fields[i].Name] < ordersMap[fields[j].Name]
	})

	return fields
}

func mapHelmPropertyTypeToFieldType(property helm.Property) string {
	switch property.Type {
	case "string":
		return "string"
	case "integer":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		return "array"
	case "object":
		if len(property.Properties) == 0 {
			return "map"
		}

		return "object"
	default:
		return property.Type
	}
}

func arrayItem(item *helm.Property) *models.Field {
	if item == nil {
		return nil
	}

	return &models.Field{
		Type:       item.Type,
		Properties: HelmSchemaToFields(*item),
	}
}

func mapTitle(name string, field helm.Property) string {
	if len(field.Title) != 0 {
		return field.Title
	}

	return name
}
