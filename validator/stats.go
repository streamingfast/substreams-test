package validator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ryanuber/columnize"
)

type Stats struct {
	entities     map[string]*entityStat
	successCount uint64
	failedCount  uint64
	totalCount   uint64
}

func (s *Stats) Print() string {
	entities := []string{}
	for ent, _ := range s.entities {
		entities = append(entities, ent)
	}
	sort.Strings(entities)
	rows := [][]string{}
	rows = append(rows, []string{"Entity", "Attr", "Total", "Success", "Failed"})

	for _, ent := range entities {
		rows = append(rows, []string{
			ent,
			"",
			fmt.Sprintf("%d", s.entities[ent].totalCount),
			ratioStr(s.entities[ent].successCount, s.entities[ent].totalCount),
			ratioStr(s.entities[ent].failedCount, s.entities[ent].totalCount),
		})

		fields := []string{}
		for fieldName, _ := range s.entities[ent].fields {
			fields = append(fields, fieldName)
		}
		sort.Strings(fields)
		for _, fieldName := range fields {
			fieldStats := s.entities[ent].fields[fieldName]
			rows = append(rows, []string{
				ent,
				fieldName,
				fmt.Sprintf("%d", fieldStats.totalCount),
				ratioStr(fieldStats.successCount, fieldStats.totalCount),
				ratioStr(fieldStats.failedCount, fieldStats.totalCount),
			})
		}
	}
	rows = append(rows, []string{"", "", fmt.Sprintf("%d", s.totalCount), ratioStr(s.successCount, s.totalCount), ratioStr(s.failedCount, s.totalCount)})

	out := []string{}
	for _, r := range rows {
		out = append(out, strings.Join(r, " | "))
	}
	return columnize.SimpleFormat(out)
}
func (s *Stats) Success(entityName, fieldName string) {
	if _, found := s.entities[entityName]; !found {
		s.entities[entityName] = &entityStat{
			entity: entityName,
			fields: map[string]*fieldStat{},
		}
	}

	if _, found := s.entities[entityName].fields[fieldName]; !found {
		s.entities[entityName].fields[fieldName] = &fieldStat{
			fieldName: fieldName,
		}
	}

	s.entities[entityName].fields[fieldName].success()
	s.entities[entityName].success()
	s.successCount++
	s.totalCount++
}

func (s *Stats) Fail(entityName, fieldName string) {
	if _, found := s.entities[entityName]; !found {
		s.entities[entityName] = &entityStat{
			entity: entityName,
			fields: map[string]*fieldStat{},
		}
	}

	if _, found := s.entities[entityName].fields[fieldName]; !found {
		s.entities[entityName].fields[fieldName] = &fieldStat{
			fieldName: fieldName,
		}
	}

	s.entities[entityName].fields[fieldName].failed()
	s.entities[entityName].failed()
	s.failedCount++
	s.totalCount++
}

func newStats() *Stats {
	return &Stats{
		entities: map[string]*entityStat{},
	}
}

type entityStat struct {
	entity       string
	fields       map[string]*fieldStat
	totalCount   uint64
	successCount uint64
	failedCount  uint64
}

func (f *entityStat) success() {
	f.totalCount++
	f.successCount++
}

func (f *entityStat) failed() {
	f.totalCount++
	f.failedCount++
}

type fieldStat struct {
	fieldName    string
	totalCount   uint64
	successCount uint64
	failedCount  uint64
}

func (f *fieldStat) success() {
	f.totalCount++
	f.successCount++
}

func (f *fieldStat) failed() {
	f.totalCount++
	f.failedCount++
}

func ratioStr(num, dem uint64) string {
	perc := float64(num) / float64(dem) * 100.0
	return fmt.Sprintf("%d (%.2f %%)", num, perc)
}
