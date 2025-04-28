cd cmd
for i in $(seq 1 100)
do
  go run main.go >> ../summary_ERCOT_6hr_random_temporal_small.log
  echo "Iteration $i/100"
  find . -name "*.log" -type f -delete
done