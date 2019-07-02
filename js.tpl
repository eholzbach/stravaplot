var encodedPoints = [
        {{ range .Poly }}
                "{{ . }}",
        {{ end }}
]
