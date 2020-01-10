url="localhost:${PORT:-8080}/jobs"

curl -s -X POST \
  -H "content-type: application/json" \
  -w " %{http_code}" \
  "$url" \
  -d @./api-examples/job.json
