package youtube

import (
	"bytes"
	"encoding/csv"
	"github.com/lithdew/bytesutil"
	"github.com/valyala/fastjson"
	"time"
)

type ListItem struct {
	ID StreamID `json:"encrypted_id"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`

	Added       string    `json:"added"`
	TimeCreated time.Time `json:"time_created"`

	Rating   float64 `json:"rating"`
	Likes    uint    `json:"likes"`
	Dislikes uint    `json:"dislikes"`

	Views    string `json:"views"`
	Comments string `json:"comments"`

	Duration      string        `json:"duration"`
	LengthSeconds time.Duration `json:"length_seconds"`

	Author  string `json:"author"`
	UserID  string `json:"user_id"`
	Privacy string `json:"privacy"`

	CategoryID uint `json:"category_id"`

	IsHD bool `json:"is_hd"`
	IsCC bool `json:"is_cc"`

	CCLicense bool `json:"cc_license"`

	Keywords []string `json:"keywords"`
}

func ParseListItem(v *fastjson.Value) ListItem {
	var r ListItem

	r.ID = StreamID(bytesutil.String(v.GetStringBytes("encrypted_id")))

	r.Title = bytesutil.String(v.GetStringBytes("title"))
	r.Description = bytesutil.String(v.GetStringBytes("description"))
	r.Thumbnail = bytesutil.String(v.GetStringBytes("thumbnail"))

	r.Added = bytesutil.String(v.GetStringBytes("added"))
	r.TimeCreated = time.Unix(v.GetInt64("time_created"), 0)

	r.Rating = v.GetFloat64("rating")
	r.Likes = v.GetUint("likes")
	r.Dislikes = v.GetUint("dislikes")

	r.Views = bytesutil.String(v.GetStringBytes("views"))
	r.Comments = bytesutil.String(v.GetStringBytes("comments"))

	r.Duration = bytesutil.String(v.GetStringBytes("duration"))
	r.LengthSeconds = time.Duration(v.GetInt64("length_seconds")) * time.Second

	r.Author = bytesutil.String(v.GetStringBytes("author"))
	r.UserID = bytesutil.String(v.GetStringBytes("user_id"))
	r.Privacy = bytesutil.String(v.GetStringBytes("privacy"))

	r.CategoryID = v.GetUint("category_id")

	r.IsHD = v.GetBool("is_hd")
	r.IsCC = v.GetBool("is_cc")

	r.CCLicense = v.GetBool("cc_license")

	fr := csv.NewReader(bytes.NewReader(v.GetStringBytes("keywords")))
	fr.Comma = ' '

	r.Keywords, _ = fr.Read()

	return r
}
