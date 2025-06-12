package models

import (
    "strings"
    "time"
)

type Date struct {
    Year  uint
    Month uint
    Day   uint
}

func NewDateFromStr(str string) (*Date, error) {
    parsed, err := time.Parse("2006-01-02", str)
    if err != nil {
        return nil, err
    }

    return NewDateFromTime(parsed), nil
}

func NewDateFromTime(time time.Time) *Date {
    newDate := Date{}
    newDate.Year = uint(time.Year())
    newDate.Month = uint(time.Month())
    newDate.Day = uint(time.Day())

    return &newDate
}

func (d *Date) UnmarshalJSON(data []byte) error {
    // Convert byte slice to string and trim double quotes
    str := strings.Trim(string(data), "\"")

    date, err := NewDateFromStr(str)
    if err != nil {
        return err
    }

    *d = *date
    return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
    return []byte(`"` + d.String() + `"`), nil
}

func (d *Date) After(other *Date) bool {
    return d.AsTime().After(other.AsTime())
}

func (d *Date) Before(other *Date) bool {
    return d.AsTime().Before(other.AsTime())
}

func (d *Date) AsTime() time.Time {
    return time.Date(
        int(d.Year),
        time.Month(int(d.Month)),
        int(d.Day),
        0,
        0,
        0,
        0,
        time.UTC,
    )
}

func (d *Date) String() string {
    return d.AsTime().Format("2006-01-02")
}

func (d *Date) Sub(other *Date) float64 {
    return d.AsTime().Sub(other.AsTime()).Hours() / 24
}
