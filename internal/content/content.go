package content

import (
	"encoding/json"
	"sync"
)

type Lesson struct {
	ID         string         `json:"id"`
	GradeID    string         `json:"gradeId"`
	LevelID    string         `json:"levelId"`
	Title      string         `json:"title"`
	Teaser     string         `json:"teaser"`
	Minutes    int            `json:"minutes"`
	Next       string         `json:"next"`
	Prev       string         `json:"prev"`
	Why        string         `json:"why,omitempty"`
	Facts      []string       `json:"facts,omitempty"`
	Scene      map[string]any `json:"scene,omitempty"`
	Trap       map[string]any `json:"trap,omitempty"`
	Board      map[string]any `json:"board"`
	Desk       map[string]any `json:"desk"`
	Brief      map[string]any `json:"brief,omitempty"`
	MaterialID string         `json:"materialId,omitempty"`
	Primer     map[string]any `json:"primer,omitempty"`
}

type Material struct {
	ID       string         `json:"id"`
	TrackID  string         `json:"trackId"`
	Title    string         `json:"title"`
	Kicker   string         `json:"kicker"`
	Grade    string         `json:"grade"`
	Free     bool           `json:"free"`
	Minutes  int            `json:"minutes"`
	LessonID string         `json:"lessonId,omitempty"`
	Teaser   string         `json:"teaser"`
	Sections map[string]any `json:"sections"`
}

type Catalog struct {
	Grades       []Grade          `json:"grades"`
	Boards       []map[string]any `json:"boards"`
	Lessons      map[string]Lesson
	Materials    []Material
	MaterialByID map[string]Material
}

type Grade struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Roman   string  `json:"roman"`
	Tagline string  `json:"tagline"`
	Levels  []Level `json:"levels"`
}

type Level struct {
	ID          string   `json:"id"`
	GradeID     string   `json:"gradeId"`
	Rank        int      `json:"rank"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	Hours       string   `json:"hours"`
	ModuleIDs   []string `json:"moduleIds"`
	MaterialIDs []string `json:"materialIds,omitempty"`
}

var (
	once sync.Once
	cat  *Catalog
	errC error
)

func Load() (*Catalog, error) {
	once.Do(func() {
		var c Catalog
		if err := json.Unmarshal(catalogJSON, &c); err != nil {
			errC = err
			return
		}
		var lessons []Lesson
		if err := json.Unmarshal(lessonsJSON, &lessons); err != nil {
			errC = err
			return
		}
		c.Lessons = map[string]Lesson{}
		for _, l := range lessons {
			c.Lessons[l.ID] = l
		}
		var boards []map[string]any
		if err := json.Unmarshal(boardsJSON, &boards); err != nil {
			errC = err
			return
		}
		c.Boards = boards
		var wrap struct {
			Topics []Material `json:"topics"`
		}
		if err := json.Unmarshal(materialsJSON, &wrap); err != nil {
			errC = err
			return
		}
		c.Materials = wrap.Topics
		c.MaterialByID = map[string]Material{}
		for _, m := range wrap.Topics {
			c.MaterialByID[m.ID] = m
		}
		cat = &c
	})
	return cat, errC
}

func (c *Catalog) Lesson(id string) (Lesson, bool) {
	l, ok := c.Lessons[id]
	return l, ok
}

func (c *Catalog) Material(id string) (Material, bool) {
	m, ok := c.MaterialByID[id]
	return m, ok
}

func (c *Catalog) Public() map[string]any {
	lessons := []map[string]any{}
	for _, g := range c.Grades {
		for _, lv := range g.Levels {
			for _, id := range lv.ModuleIDs {
				l, ok := c.Lessons[id]
				if !ok {
					continue
				}
				lessons = append(lessons, map[string]any{
					"id": l.ID, "gradeId": l.GradeID, "levelId": l.LevelID,
					"title": l.Title, "teaser": l.Teaser, "minutes": l.Minutes,
					"next": l.Next, "prev": l.Prev,
				})
			}
		}
	}
	return map[string]any{"grades": c.Grades, "lessons": lessons, "boards": boardMeta(c.Boards)}
}

func boardMeta(boards []map[string]any) []map[string]any {
	out := []map[string]any{}
	for _, b := range boards {
		out = append(out, map[string]any{
			"id": b["id"], "title": b["title"], "teaser": b["teaser"],
			"term": b["term"], "learn": b["learn"],
			"grade": b["grade"], "plan": b["plan"], "minutes": b["minutes"],
		})
	}
	return out
}
