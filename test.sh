cd cmd

for i in $(seq 62 100)
do
  go run main.go MISO 6hr random hybridSelection 0.7 >> ../summary_MISO_6hr_random_hybridSelection_70.log
  echo "Iteration $i/100"
  find . -name "*.log" -type f -delete
done

for i in $(seq 1 100)
do
  go run main.go MISO 6hr random hybridSelection 0.8 >> ../summary_MISO_6hr_random_hybridSelection_80.log
  echo "Iteration $i/100"
  find . -name "*.log" -type f -delete
done

for i in $(seq 1 100)
do
  go run main.go MISO 6hr random hybridSelection 0.9 >> ../summary_MISO_6hr_random_hybridSelection_90.log
  echo "Iteration $i/100"
  find . -name "*.log" -type f -delete
done