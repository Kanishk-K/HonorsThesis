cd cmd
for i in $(seq 1 5)
do
  go run main.go >> ../summary_6hr_random_temporal_small.log
done