{{- define "cilikube.name" -}}
{{- default .Chart.Name .Values.app.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cilikube.fullname" -}}
{{- printf "%s" (include "cilikube.name" .) -}}
{{- end -}}

{{- define "cilikube.labels" -}}
app.kubernetes.io/name: {{ include "cilikube.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}
