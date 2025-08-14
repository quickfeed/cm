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

| Lab | Topic                                                     | Grading          | Approval             | Submission              | Deadline          |
|:---:|-----------------------------------------------------------|------------------|----------------------|-------------------------|-------------------|
{{- range $index, $a := .}}
| {{ $a.Order }} | [{{ $a.Title }}][{{ $a.Order }}] | {{ $a.Grading }} | {{ $a.ApproveType }} | {{ $a.SubmissionType }} | {{ $a.ShortDeadline }} |
{{- end}}
{{range $index, $a := .}}
[{{ $a.Order }}]: https://github.com/{{$course}}-{{$year}}/assignments/tree/main/{{ $a.Name }}
{{- end}}
`
