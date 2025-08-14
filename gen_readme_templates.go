package main

const titleOnlyHeaderTemplate = `# Lab {{ .Order }}: {{ .Title }}
`

const headerTemplate = `# Lab {{ .Order }}: {{ .Title }}

| Lab {{ .Order }}: | {{ .Title }} |
| ---------------------    | --------------------- |
| Subject:                 | {{ .Subject }} |
| Deadline:                | **{{ .Deadline }}** |
| Expected effort:         | {{ .HoursMin -}}-{{- .HoursMax}} hours |
| Grading:                 | {{ .Grading }} |
| Submission:              | {{ .SubmissionType }} |
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
