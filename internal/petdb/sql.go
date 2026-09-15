package petdb

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func trimSQL(sqlText string) string {
	s := strings.TrimSpace(stripSQLComments(sqlText))
	s = strings.TrimRight(s, "; \t\n\r")
	return strings.TrimSpace(s)
}

func execSelect(sqlText string, tables map[string][]map[string]any) ([]string, []map[string]any, error) {
	s := trimSQL(sqlText)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	low := strings.ToLower(s)
	if !strings.HasPrefix(low, "select") {
		return nil, nil, fmt.Errorf("ожидался SELECT")
	}
	body := strings.TrimSpace(s[6:])
	fromIdx := indexKW(body, " from ")
	if fromIdx < 0 {
		return nil, nil, fmt.Errorf("нужен FROM")
	}
	selectList := strings.TrimSpace(body[:fromIdx])
	rest := strings.TrimSpace(body[fromIdx+6:])
	whereIdx := indexKW(rest, " where ")
	groupIdx := indexKW(rest, " group by ")
	havingIdx := indexKW(rest, " having ")
	orderIdx := indexKW(rest, " order by ")

	cut := len(rest)
	for _, i := range []int{whereIdx, groupIdx, havingIdx, orderIdx} {
		if i >= 0 && i < cut {
			cut = i
		}
	}
	fromPart := strings.TrimSpace(rest[:cut])
	var where, group, having, order string
	if whereIdx >= 0 {
		end := len(rest)
		for _, i := range []int{groupIdx, havingIdx, orderIdx} {
			if i > whereIdx && i < end {
				end = i
			}
		}
		where = strings.TrimSpace(rest[whereIdx+7 : end])
	}
	if groupIdx >= 0 {
		end := len(rest)
		for _, i := range []int{havingIdx, orderIdx} {
			if i > groupIdx && i < end {
				end = i
			}
		}
		group = strings.TrimSpace(rest[groupIdx+10 : end])
	}
	if havingIdx >= 0 {
		end := len(rest)
		if orderIdx > havingIdx {
			end = orderIdx
		}
		having = strings.TrimSpace(rest[havingIdx+8 : end])
	}
	if orderIdx >= 0 {
		order = strings.TrimSpace(rest[orderIdx+10:])
	}

	rows, err := fromJoin(fromPart, tables)
	if err != nil {
		return nil, nil, err
	}
	if where != "" {
		filtered := []map[string]any{}
		for _, r := range rows {
			ok, err := evalBool(where, r)
			if err != nil {
				return nil, nil, err
			}
			if ok {
				filtered = append(filtered, r)
			}
		}
		rows = filtered
	}

	aliases, agg := parseSelectList(selectList)
	if group != "" || agg {
		rows, err = groupBy(rows, group, aliases)
		if err != nil {
			return nil, nil, err
		}
		if having != "" {
			filtered := []map[string]any{}
			for _, r := range rows {
				ok, err := evalBool(having, r)
				if err != nil {
					return nil, nil, err
				}
				if ok {
					filtered = append(filtered, r)
				}
			}
			rows = filtered
		}
	} else {
		projected := []map[string]any{}
		for _, r := range rows {
			p, err := projectRow(r, aliases)
			if err != nil {
				return nil, nil, err
			}
			projected = append(projected, p)
		}
		rows = projected
	}
	if order != "" {
		rows = orderBy(rows, order)
	}
	cols := make([]string, 0, len(aliases))
	for _, a := range aliases {
		cols = append(cols, a.name)
	}
	if len(cols) == 0 && len(rows) > 0 {
		for k := range rows[0] {
			cols = append(cols, k)
		}
	}
	return cols, rows, nil
}

type selCol struct {
	expr string
	name string
	agg  bool
}

func parseSelectList(s string) ([]selCol, bool) {
	parts := splitComma(s)
	out := []selCol{}
	agg := false
	for _, p := range parts {
		p = strings.TrimSpace(p)
		name := ""
		expr := p
		if i := indexKW(p, " as "); i >= 0 {
			expr = strings.TrimSpace(p[:i])
			name = strings.TrimSpace(p[i+4:])
		} else {
			fs := strings.Fields(p)
			if len(fs) >= 2 && !strings.ContainsAny(fs[len(fs)-1], "()+-*/") {
				maybe := fs[len(fs)-1]
				if !isKW(maybe) && !strings.Contains(maybe, ".") {
					name = maybe
					expr = strings.TrimSpace(strings.Join(fs[:len(fs)-1], " "))
				}
			}
		}
		name = strings.Trim(name, "`\"")
		if name == "" {
			name = strings.ToLower(strings.ReplaceAll(expr, " ", "_"))
			if i := strings.LastIndex(name, "."); i >= 0 {
				name = name[i+1:]
			}
			name = regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(name, "_")
			name = strings.Trim(name, "_")
			if name == "" {
				name = "col"
			}
		}
		isAgg := strings.Contains(strings.ToLower(expr), "sum(") || strings.Contains(strings.ToLower(expr), "count(")
		if isAgg {
			agg = true
		}
		out = append(out, selCol{expr: expr, name: name, agg: isAgg})
	}
	return out, agg
}

func fromJoin(fromPart string, tables map[string][]map[string]any) ([]map[string]any, error) {
	reJoin := regexp.MustCompile(`(?i)\s+(?:inner\s+)?join\s+`)
	chunks := reJoin.Split(fromPart, -1)
	type src struct {
		table, alias, on string
	}
	parts := []src{}
	first := strings.Fields(chunks[0])
	if len(first) == 0 {
		return nil, fmt.Errorf("пустой FROM")
	}
	t0 := strings.Trim(strings.ToLower(first[0]), ",;`\"")
	a0 := t0
	if len(first) >= 2 && strings.ToLower(first[1]) == "as" && len(first) >= 3 {
		a0 = strings.ToLower(first[2])
	} else if len(first) >= 2 && !isKW(first[1]) {
		a0 = strings.ToLower(first[1])
	}
	parts = append(parts, src{table: t0, alias: a0})
	for i := 1; i < len(chunks); i++ {
		ch := strings.TrimSpace(chunks[i])
		onIdx := indexKW(ch, " on ")
		if onIdx < 0 {
			return nil, fmt.Errorf("JOIN без ON")
		}
		head := strings.TrimSpace(ch[:onIdx])
		on := strings.TrimSpace(ch[onIdx+4:])
		fs := strings.Fields(head)
		tb := strings.Trim(strings.ToLower(fs[0]), ",;`\"")
		al := tb
		if len(fs) >= 3 && strings.ToLower(fs[1]) == "as" {
			al = strings.ToLower(fs[2])
		} else if len(fs) >= 2 {
			al = strings.ToLower(fs[1])
		}
		parts = append(parts, src{table: tb, alias: al, on: on})
	}
	base, ok := tables[parts[0].table]
	if !ok {
		return nil, fmt.Errorf("нет таблицы %s", parts[0].table)
	}
	rows := prefixRows(base, parts[0].alias)
	for _, p := range parts[1:] {
		right, ok := tables[p.table]
		if !ok {
			return nil, fmt.Errorf("нет таблицы %s", p.table)
		}
		rr := prefixRows(right, p.alias)
		joined := []map[string]any{}
		for _, L := range rows {
			for _, R := range rr {
				m := mergeRow(L, R)
				ok, err := evalBool(p.on, m)
				if err != nil {
					return nil, err
				}
				if ok {
					joined = append(joined, m)
				}
			}
		}
		rows = joined
	}
	return rows, nil
}

func prefixRows(rows []map[string]any, alias string) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		m := map[string]any{}
		for k, v := range r {
			m[k] = v
			m[alias+"."+k] = v
		}
		out = append(out, m)
	}
	return out
}

func mergeRow(a, b map[string]any) map[string]any {
	m := map[string]any{}
	for k, v := range a {
		m[k] = v
	}
	for k, v := range b {
		m[k] = v
	}
	return m
}

func projectRow(r map[string]any, cols []selCol) (map[string]any, error) {
	out := map[string]any{}
	for _, c := range cols {
		if strings.TrimSpace(c.expr) == "*" {
			for k, v := range r {
				if !strings.Contains(k, ".") {
					out[k] = v
				}
			}
			continue
		}
		v, err := evalExpr(c.expr, r)
		if err != nil {
			return nil, err
		}
		out[c.name] = v
	}
	return out, nil
}

func groupBy(rows []map[string]any, group string, cols []selCol) ([]map[string]any, error) {
	keys := []string{}
	if group != "" {
		for _, g := range splitComma(group) {
			keys = append(keys, strings.TrimSpace(g))
		}
	}
	buckets := [][]map[string]any{}
	index := map[string]int{}
	for _, r := range rows {
		kb := ""
		for _, k := range keys {
			v, _ := evalExpr(k, r)
			kb += fmt.Sprintf("%v|", v)
		}
		if keys == nil || group == "" {
			kb = "*"
		}
		i, ok := index[kb]
		if !ok {
			i = len(buckets)
			index[kb] = i
			buckets = append(buckets, nil)
		}
		buckets[i] = append(buckets[i], r)
	}
	out := []map[string]any{}
	for _, b := range buckets {
		row := map[string]any{}
		for _, c := range cols {
			if c.agg {
				v, err := evalAgg(c.expr, b)
				if err != nil {
					return nil, err
				}
				row[c.name] = v
				continue
			}
			v, err := evalExpr(c.expr, b[0])
			if err != nil {
				return nil, err
			}
			row[c.name] = v
		}
		row["count(*)"] = float64(len(b))
		out = append(out, row)
	}
	return out, nil
}

func evalAgg(expr string, rows []map[string]any) (any, error) {
	e := strings.TrimSpace(expr)
	low := strings.ToLower(e)
	if strings.HasPrefix(low, "count(") {
		inner := e[6 : len(e)-1]
		if strings.TrimSpace(inner) == "*" {
			return float64(len(rows)), nil
		}
		n := 0
		for _, r := range rows {
			v, _ := evalExpr(inner, r)
			if v != nil {
				n++
			}
		}
		return float64(n), nil
	}
	if strings.HasPrefix(low, "sum(") {
		inner := e[4 : len(e)-1]
		sum := 0.0
		for _, r := range rows {
			v, err := evalExpr(inner, r)
			if err != nil {
				return nil, err
			}
			sum += num(v)
		}
		return sum, nil
	}
	return evalExpr(expr, rows[0])
}

func orderBy(rows []map[string]any, order string) []map[string]any {
	key := strings.TrimSpace(strings.Split(order, " ")[0])
	desc := strings.Contains(strings.ToLower(order), " desc")
	cp := append([]map[string]any{}, rows...)
	for i := 0; i < len(cp); i++ {
		for j := i + 1; j < len(cp); j++ {
			vi, _ := evalExpr(key, cp[i])
			vj, _ := evalExpr(key, cp[j])
			less := fmt.Sprint(vi) < fmt.Sprint(vj)
			if numish(vi) || numish(vj) {
				less = num(vi) < num(vj)
			}
			if desc {
				less = !less
			}
			if !less {
				cp[i], cp[j] = cp[j], cp[i]
			}
		}
	}
	return cp
}

func numish(v any) bool {
	switch v.(type) {
	case float64, float32, int, int64:
		return true
	}
	return false
}

func evalBool(expr string, row map[string]any) (bool, error) {
	v, err := evalExpr(expr, row)
	if err != nil {
		return false, err
	}
	switch t := v.(type) {
	case bool:
		return t, nil
	default:
		return num(t) != 0 && t != nil, nil
	}
}

func evalExpr(expr string, row map[string]any) (any, error) {
	e := strings.TrimSpace(expr)
	if e == "" {
		return nil, nil
	}
	low := strings.ToLower(e)
	if i := indexKW(e, " or "); i >= 0 {
		l, err := evalBool(e[:i], row)
		if err != nil {
			return nil, err
		}
		r, err := evalBool(e[i+4:], row)
		if err != nil {
			return nil, err
		}
		return l || r, nil
	}
	if i := indexKW(e, " and "); i >= 0 {
		l, err := evalBool(e[:i], row)
		if err != nil {
			return nil, err
		}
		r, err := evalBool(e[i+5:], row)
		if err != nil {
			return nil, err
		}
		return l && r, nil
	}
	if strings.HasPrefix(low, "not ") {
		v, err := evalBool(e[4:], row)
		return !v, err
	}
	if strings.HasPrefix(e, "(") && strings.HasSuffix(e, ")") && balanced(e[1:len(e)-1]) {
		return evalExpr(e[1:len(e)-1], row)
	}
	if strings.HasPrefix(low, "case ") || strings.HasPrefix(low, "case\n") || strings.HasPrefix(low, "case\t") {
		return evalCase(e, row)
	}
	if strings.HasPrefix(low, "datetime(") {
		return evalDatetime(e)
	}
	if strings.HasPrefix(low, "substr(") {
		return evalSubstr(e, row)
	}
	if strings.HasPrefix(low, "coalesce(") {
		inner := e[9 : len(e)-1]
		for _, p := range splitComma(inner) {
			v, err := evalExpr(p, row)
			if err != nil {
				return nil, err
			}
			if v != nil && str(v) != "" {
				return v, nil
			}
		}
		return nil, nil
	}
	for _, op := range []string{">=", "<=", "!=", "<>", "=", ">", "<"} {
		if i := indexTop(e, op); i >= 0 {
			l, err := evalExpr(strings.TrimSpace(e[:i]), row)
			if err != nil {
				return nil, err
			}
			r, err := evalExpr(strings.TrimSpace(e[i+len(op):]), row)
			if err != nil {
				return nil, err
			}
			return cmp(l, r, op), nil
		}
	}
	if i := indexKW(e, " is not null"); i >= 0 {
		v, err := evalExpr(e[:i], row)
		return v != nil && str(v) != "", err
	}
	if i := indexKW(e, " is null"); i >= 0 {
		v, err := evalExpr(e[:i], row)
		return v == nil || str(v) == "", err
	}
	if i := indexKW(e, " in "); i >= 0 {
		left, err := evalExpr(e[:i], row)
		if err != nil {
			return nil, err
		}
		list := strings.TrimSpace(e[i+4:])
		list = strings.TrimPrefix(list, "(")
		list = strings.TrimSuffix(list, ")")
		for _, p := range splitComma(list) {
			rv, err := evalExpr(p, row)
			if err != nil {
				return nil, err
			}
			if fmt.Sprint(left) == fmt.Sprint(rv) {
				return true, nil
			}
		}
		return false, nil
	}
	if i := indexTop(e, "+"); i > 0 {
		l, _ := evalExpr(e[:i], row)
		r, _ := evalExpr(e[i+1:], row)
		return num(l) + num(r), nil
	}
	if i := indexTop(e, "-"); i > 0 {
		l, _ := evalExpr(e[:i], row)
		r, _ := evalExpr(e[i+1:], row)
		return num(l) - num(r), nil
	}
	if i := indexTop(e, "*"); i > 0 {
		l, _ := evalExpr(e[:i], row)
		r, _ := evalExpr(e[i+1:], row)
		return num(l) * num(r), nil
	}
	if i := indexTop(e, "/"); i > 0 {
		l, _ := evalExpr(e[:i], row)
		r, _ := evalExpr(e[i+1:], row)
		if num(r) == 0 {
			return 0.0, nil
		}
		return num(l) / num(r), nil
	}
	if (strings.HasPrefix(e, "'") && strings.HasSuffix(e, "'")) || (strings.HasPrefix(e, `"`) && strings.HasSuffix(e, `"`)) {
		return e[1 : len(e)-1], nil
	}
	if low == "null" {
		return nil, nil
	}
	if n, err := strconv.ParseFloat(e, 64); err == nil {
		return n, nil
	}
	if v, ok := lookup(row, e); ok {
		return v, nil
	}
	return nil, fmt.Errorf("не понял выражение: %s", e)
}

func evalCase(e string, row map[string]any) (any, error) {
	s := e
	low := strings.ToLower(s)
	if strings.HasSuffix(low, " end") {
		s = strings.TrimSpace(s[:len(s)-4])
	}
	s = strings.TrimSpace(s[4:])
	elseVal := any(nil)
	if i := indexKW(s, " else "); i >= 0 {
		ev, err := evalExpr(s[i+6:], row)
		if err != nil {
			return nil, err
		}
		elseVal = ev
		s = strings.TrimSpace(s[:i])
	}
	for {
		wi := indexKW(s, "when ")
		if wi < 0 {
			break
		}
		s = strings.TrimSpace(s[wi+5:])
		ti := indexKW(s, " then ")
		if ti < 0 {
			break
		}
		cond := strings.TrimSpace(s[:ti])
		s = strings.TrimSpace(s[ti+6:])
		next := indexKW(s, " when ")
		thenExpr := s
		if next >= 0 {
			thenExpr = strings.TrimSpace(s[:next])
			s = s[next:]
		} else {
			s = ""
		}
		ok, err := evalBool(cond, row)
		if err != nil {
			return nil, err
		}
		if ok {
			return evalExpr(thenExpr, row)
		}
		if s == "" {
			break
		}
	}
	return elseVal, nil
}

func evalDatetime(e string) (any, error) {
	inner := e[9 : len(e)-1]
	parts := splitComma(inner)
	base := time.Now().UTC()
	if len(parts) > 0 && strings.Trim(strings.ToLower(parts[0]), "'\"") == "now" {
		base = time.Now().UTC()
	}
	if len(parts) > 1 {
		mod := strings.Trim(parts[1], " '\"")
		if strings.HasPrefix(mod, "-") && strings.Contains(mod, "hour") {
			n := 12
			fmt.Sscanf(mod, "-%d", &n)
			base = base.Add(-time.Duration(n) * time.Hour)
		}
	}
	return base.Format("2006-01-02T15:04:05.000Z"), nil
}

func evalSubstr(e string, row map[string]any) (any, error) {
	inner := e[7 : len(e)-1]
	ps := splitComma(inner)
	if len(ps) < 3 {
		return "", fmt.Errorf("substr(s, start, len)")
	}
	v, _ := evalExpr(ps[0], row)
	start := int(num(mustNum(ps[1])))
	ln := int(num(mustNum(ps[2])))
	s := str(v)
	if start < 1 {
		start = 1
	}
	i := start - 1
	if i >= len(s) {
		return "", nil
	}
	j := i + ln
	if j > len(s) {
		j = len(s)
	}
	return s[i:j], nil
}

func mustNum(s string) any {
	n, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return n
}

func lookup(row map[string]any, name string) (any, bool) {
	n := strings.Trim(name, "`\"")
	if v, ok := row[n]; ok {
		return v, true
	}
	low := strings.ToLower(n)
	for k, v := range row {
		if strings.ToLower(k) == low {
			return v, true
		}
		if strings.HasSuffix(strings.ToLower(k), "."+low) {
			return v, true
		}
	}
	if i := strings.LastIndex(n, "."); i >= 0 {
		return lookup(row, n[i+1:])
	}
	return nil, false
}

func cmp(l, r any, op string) bool {
	if op == "=" {
		if l == nil || r == nil {
			return l == nil && r == nil
		}
		if numish(l) || numish(r) {
			return num(l) == num(r)
		}
		return fmt.Sprint(l) == fmt.Sprint(r)
	}
	if op == "!=" || op == "<>" {
		return !cmp(l, r, "=")
	}
	ln, rn := fmt.Sprint(l), fmt.Sprint(r)
	if numish(l) || numish(r) || isTimeish(ln) {
		if isTimeish(ln) || isTimeish(rn) {
			return cmpTime(ln, rn, op)
		}
		a, b := num(l), num(r)
		switch op {
		case ">":
			return a > b
		case "<":
			return a < b
		case ">=":
			return a >= b
		case "<=":
			return a <= b
		}
	}
	switch op {
	case ">":
		return ln > rn
	case "<":
		return ln < rn
	case ">=":
		return ln >= rn
	case "<=":
		return ln <= rn
	}
	return false
}

func isTimeish(s string) bool {
	return len(s) >= 10 && (strings.Contains(s, "T") || strings.Contains(s, "-")) && s[4] == '-'
}

func cmpTime(a, b, op string) bool {
	switch op {
	case ">":
		return a > b
	case "<":
		return a < b
	case ">=":
		return a >= b
	case "<=":
		return a <= b
	}
	return false
}

func indexKW(s, kw string) int {
	low := strings.ToLower(s)
	k := strings.ToLower(kw)
	depth := 0
	for i := 0; i <= len(low)-len(k); i++ {
		c := low[i]
		if c == '(' {
			depth++
		} else if c == ')' && depth > 0 {
			depth--
		}
		if depth == 0 && low[i:i+len(k)] == k {
			return i
		}
	}
	return -1
}

func indexTop(s, op string) int {
	depth := 0
	for i := 0; i <= len(s)-len(op); i++ {
		if s[i] == '(' {
			depth++
			continue
		}
		if s[i] == ')' && depth > 0 {
			depth--
			continue
		}
		if depth == 0 && s[i:i+len(op)] == op {
			if op == "-" && i == 0 {
				continue
			}
			return i
		}
	}
	return -1
}

func splitComma(s string) []string {
	out := []string{}
	depth := 0
	start := 0
	for i, c := range s {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	out = append(out, strings.TrimSpace(s[start:]))
	return out
}

func balanced(s string) bool {
	d := 0
	for _, c := range s {
		if c == '(' {
			d++
		}
		if c == ')' {
			d--
			if d < 0 {
				return false
			}
		}
	}
	return d == 0
}

func isKW(s string) bool {
	switch strings.ToLower(s) {
	case "on", "as", "from", "where", "join", "inner", "left", "group", "by", "having", "order", "and", "or", "not", "in":
		return true
	}
	return false
}
