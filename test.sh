cd cmd
for i in $(seq 1 100)
do
  go run main.go >> ../summary_CAISO_6hr_random_modelSelection_small.log
  find . -name "*.log" -type f -delete
done