package Utils

import (
	"encoding/json"
	"log"
	"os"
)

// Config : TCP Conf
type Config struct {
	TCP struct {
		Host string
		Port int
	}
}

// type Config struct {
//     TCP []struct {
//         Host string
//         Port int
//     }
// }
// === FOR THIS json ===
// {
//     "tcp": [{
//         "host": "xxx",
//         "port": xxx
//     },{
//         "host": "xxx",
//         "port": xxx
//     }]
// }

func readTCPConf(filePath string) (c Config) {
	r, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Read TCP Conf Error --> ", err)
	}

	decoder := json.NewDecoder(r)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatalln("Decode TCP Conf Error --> ", err)
	}
	return c
}

// GetTCPConf return TCP Config
// @param  filePath string
// @return host string
// @return port int
func GetTCPConf(filePath string) (host string, port int) {
	c := readTCPConf(filePath)
	return c.TCP.Host, c.TCP.Port
}
