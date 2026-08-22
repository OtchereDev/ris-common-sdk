package export

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const MaxRows = 50000

type Column struct{ Key, Header string }

func ResolveColumns(resource, requested string) ([]Column, error) {
	if strings.TrimSpace(requested) == "" {
		return nil, errors.New("unknown resource: empty resouce")
	}

	all, ok := Columns(resource)
	if !ok {
		return nil, fmt.Errorf("unknown resource %q", resource)
	}

	allowed := make(map[string]struct{}, len(all))
	for _, c := range all {
		allowed[c.Key] = struct{}{}
	}

	wanted := make(map[string]struct{}, len(all))
	for key := range strings.SplitSeq(requested, ",") {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, ok := allowed[key]; !ok {
			return nil, fmt.Errorf("unknown column %q for resource %q", key, resource)
		}
		wanted[key] = struct{}{}
	}

	if len(wanted) == 0 {
		return all, nil
	}

	resolved := make([]Column, 0, len(wanted))
	for _, c := range all {
		if _, ok := wanted[c.Key]; ok {
			resolved = append(resolved, c)
		}
	}
	return resolved, nil
}

var acronyms = map[string]string{
	"id":  "ID",
	"dob": "DOB",
	"sms": "SMS",
	"url": "URL",
}

func Columns(resource string) ([]Column, bool) {
	keys, ok := mapping[resource]
	if !ok {
		return nil, false
	}

	cols := make([]Column, 0, len(keys))
	for _, key := range keys {
		cols = append(cols, Column{Key: key, Header: header(key)})
	}
	return cols, true
}

func header(key string) string {
	parts := strings.Split(key, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		if a, ok := acronyms[part]; ok {
			parts[i] = a
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}

// WriteCSV streams rows in the given column order.
func WriteCSV(w io.Writer, cols []Column, rows []map[string]any) error {
	if len(cols) == 0 {
		return errors.New("no columns to write")
	}
	if len(rows) > MaxRows {
		return fmt.Errorf("too many rows: %d exceeds limit of %d", len(rows), MaxRows)
	}

	buf := bufio.NewWriterSize(w, 64<<10)
	cw := csv.NewWriter(buf)

	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.Header
	}
	if err := cw.Write(headers); err != nil {
		return err
	}

	record := make([]string, len(cols))
	for _, row := range rows {
		for i, c := range cols {
			record[i] = cell(row[c.Key])
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}

	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}
	return buf.Flush()
}

func cell(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case []byte:
		return string(t)
	case time.Time:
		if t.IsZero() {
			return ""
		}
		return t.Format(time.RFC3339)
	case *time.Time:
		if t == nil || t.IsZero() {
			return ""
		}
		return t.Format(time.RFC3339)
	case bool:
		if t {
			return "Yes"
		}
		return "No"
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	default:
		return fmt.Sprint(t)
	}
}

// Attach sets Content-Type and Content-Disposition with a dated filename.
func Attach(c *fiber.Ctx, basename string) {
	name := safeName(basename)
	if name == "" {
		name = "export"
	}

	filename := fmt.Sprintf("%s-%s.csv", name, time.Now().Format("2006-01-02"))
	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", filename))
}

func safeName(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 4 && strings.EqualFold(s[len(s)-4:], ".csv") {
		s = s[:len(s)-4]
	}

	s = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			return r
		case r == '-' || r == '_':
			return r
		case r == ' ' || r == '.' || r == '/' || r == '\\':
			return '-'
		default:
			return -1
		}
	}, s)

	return strings.Trim(s, "-")
}
