package main

import (
	"encoding/csv"
	"github.com/eranet/rhombus/rhomgo"
	"log"
	"os"
	"strconv"
)

type Position struct {
	Value float64
}

type Encoders struct {
	Pos1 float64
	Vel1 float64
	Pos2 float64
	Vel2 float64
}

var CommandSubject = "/pendabot/shoulder_torque_controller/command"

func main() {
	c := rhomgo.LocalJSONConnection()
	defer c.Close()

	latestEncoders := Encoders{}

	c.Subscribe("/encoders", func(p *Encoders) {
		//fmt.Printf("Received a position: %+v\n", p)
		latestEncoders = *p
	})

	rate := rhomgo.NewRate(1)
	for i := 1.0; i > -1; i -= 0.001 {
		err := c.Publish(CommandSubject, Position{Value: i})
		if err != nil {
			println("error pub:", err)
		}
		rate.Sleep()
	}

	N := 200
	traj := readTrajectoryCsvFile("trajectory.csv")
	policy := readPolicyCsvFile("policy.csv")

	//data = np.zeros((N,8), dtype=float)
	//Kdu := []float64{16.4736, 35.6838, 0.0973, 4.0886}

	for i := 0; i < N; i++ {
		//t = time.time()
		X := []float64{latestEncoders.Pos1, latestEncoders.Pos2, latestEncoders.Vel1, latestEncoders.Vel2}
		X_des := traj[i][0:4]
		u_des := traj[i][4]
		K := policy[i]
		//u := u_des - 1.0*np.dot(K, X_des-X)
		u := u_des - (K[0]*(X_des[0]-X[0]) + K[1]*(X_des[1]-X[1]) + K[2]*(X_des[2]-X[2]) + K[3]*(X_des[3]-X[3]))

		if u > 10 { //safety to not pass big value of current
			u = 10
		}
		if u < -10 {
			u = -10
		}

		err := c.Publish(CommandSubject, Position{Value: u})
		if err != nil {
			println("error pub:", err)
		}
		//u_meas = -my_drive.axis1.motor.current_control.Iq_measured
	}

}
func readTrajectoryCsvFile(filePath string) [][]float64 {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	res := [][]float64{}
	for i := 0; i < len(data[0]); i++ {
		res = append(res, []float64{toF(data[0][i]), toF(data[1][i]), toF(data[2][i]), toF(data[3][i]), toF(data[4][i])})
	}
	return res
}

func readPolicyCsvFile(filePath string) [][]float64 {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	res := [][]float64{}
	for i := 0; i < len(data); i++ {
		res = append(res, []float64{toF(data[i][0]), toF(data[i][1]), toF(data[i][2]), toF(data[i][3])})
	}
	return res
}

func toF(s string) float64 {
	res, _ := strconv.ParseFloat(s, 64)
	return res
}
