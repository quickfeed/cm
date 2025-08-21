package main

const titleOnlyHeaderTemplate = `# Lab {{ .Order }}: {{ .Title }}
`

const headerTemplate = `# Lab {{ .Order }}: {{ .Title }}
{{- $col1Width := col1Width . }}
{{- $col2Width := col2Width . }}
{{- $labels := headerLabels . }}
{{- $values := headerValues . }}
{{range $index, $label := $labels }}
{{- $value := index $values $index }}
| {{ padRight $label $col1Width }} | {{ padRight $value $col2Width }} |
{{- if eq $index 0 }}
| {{ separator $col1Width }} | {{ separator $col2Width }} |
{{- end }}
{{- end }}
`

const tocTemplate = `{{- .LabHeader }}
## Table of Contents

- [Table of Contents](#table-of-contents)
{{- range $index, $heading := slice .ToC 1}}
{{- if hasPrefix $heading "## "}}
- [{{escapeText (trimPrefix $heading "## ")}}](#{{link (trimPrefix $heading "## ")}})
{{- else if hasPrefix $heading "### "}}
  - [{{escapeText (trimPrefix $heading "### ")}}](#{{link (trimPrefix $heading "### ")}})
{{- else}}
- [{{escapeText $heading}}](#{{link $heading}})
{{- end}}
{{- end}}

`

const labPlanTemplate = `{{- $assignment := index . 1}}
{{- $year := $assignment.Year}}
{{- $course := $assignment.CourseOrg -}}
# Lab Plan for {{$year}}

{{- $labs := . }}
{{- $col1Width := len " Lab " }}
{{- $col2Width := add (len "Topic") 5 }}
{{- $col3Width := len "Grading" }}
{{- $col4Width := len "Approval" }}
{{- $col5Width := len "Submission" }}
{{- $col6Width := len "Deadline" }}
{{- range $i, $a := $labs }}
  {{- if gt (len (printf "%v" $a.Order)) $col1Width }}{{ $col1Width = len (printf "%v" $a.Order) }}{{- end }}
  {{- if gt (add (len $a.Title) 5) $col2Width }}{{ $col2Width = add (len $a.Title) 5 }}{{- end }}
  {{- if gt (len $a.Grading) $col3Width }}{{ $col3Width = len $a.Grading }}{{- end }}
  {{- if gt (len $a.ApproveType) $col4Width }}{{ $col4Width = len $a.ApproveType }}{{- end }}
  {{- if gt (len $a.SubmissionType) $col5Width }}{{ $col5Width = len $a.SubmissionType }}{{- end }}
  {{- if gt (len $a.ShortDeadline) $col6Width }}{{ $col6Width = len $a.ShortDeadline }}{{- end }}
{{- end }}

| {{ padRight " Lab " $col1Width }} | {{ padRight "Topic" $col2Width }} | {{ padRight "Grading" $col3Width }} | {{ padRight "Approval" $col4Width }} | {{ padRight "Submission" $col5Width }} | {{ padRight "Deadline" $col6Width }} |
| {{ padRight ":---:" $col1Width }} | {{ padRight (separator $col2Width) $col2Width }} | {{ padRight (separator $col3Width) $col3Width }} | {{ padRight (separator $col4Width) $col4Width }} | {{ padRight (separator $col5Width) $col5Width }} | {{ padRight (separator $col6Width) $col6Width }} |
{{- range $index, $a := .}}
| {{ center (printf "%v" $a.Order) $col1Width }} | {{ padRight (printf "[%s][%v]" $a.Title $a.Order) $col2Width }} | {{ padRight $a.Grading $col3Width }} | {{ padRight $a.ApproveType $col4Width }} | {{ padRight $a.SubmissionType $col5Width }} | {{ padRight $a.ShortDeadline $col6Width }} |
{{- end}}
{{range $index, $a := .}}
[{{ $a.Order }}]: https://github.com/{{$course}}-{{$year}}/assignments/tree/main/{{ $a.Name }}
{{- end}}
`
