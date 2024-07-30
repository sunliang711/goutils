package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sunliang711/goutils/rmq"
)

func main() {

	caCertBytes := []byte(`-----BEGIN CERTIFICATE-----
MIIDtjCCAp6gAwIBAgITCeBKIs7Ew4XSsy/Wul3kCQ9KsTANBgkqhkiG9w0BAQsF
ADBrMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwN
U2FuIEZyYW5jaXNjbzEOMAwGA1UECgwFTXlPcmcxDzANBgNVBAsMBk15RGVwdDEO
MAwGA1UEAwwFbXktY2EwHhcNMjQwNzIyMTA0MTAzWhcNMjUwNzIyMTA0MTAzWjBr
MQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2Fu
IEZyYW5jaXNjbzEOMAwGA1UECgwFTXlPcmcxDzANBgNVBAsMBk15RGVwdDEOMAwG
A1UEAwwFbXktY2EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDXiqz/
sS/KeoL4nVmgZjYNO9mZpZz2dVA5C5Ep8E7HFF93/ctrsdukZ+MMi5u4rfhDO6xn
dL9FV7A/IwYWBNXheuN5DDhSIyAu8ZTRHzEFYpQ4l2zLqIb04FaONCO1LReh30SX
qv97n7BihiqqVnpQLpriScKiVk4Cw21LyH6BxkzXqgMl0HmvJ8mg1IQlqkTQhPHn
ctsB91uqyiShPU841bqskkpwSXXdLOVbBIQ8IMumTCOoWdlckvqcnfZsrVcLGYGa
vNuaUbEFreMNpYN8r8nVMJCQIpUdqN8I7KFDAxOqgI/StwNFrWVEdyAlMUpMJKZ4
IGPAaRkAXxr7aEJZAgMBAAGjUzBRMB0GA1UdDgQWBBSq/ym/7JkJW/+BGwHFsEch
9+oTsjAfBgNVHSMEGDAWgBSq/ym/7JkJW/+BGwHFsEch9+oTsjAPBgNVHRMBAf8E
BTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQDDXdFxJJBKuCIYoCyY/NTR7lok3Hk6
CvT97LMztoiIciR7MJqNK7e1xgFFs5DfSgYqdzSdUjvUG+WoUk/jXJzSsanKsErF
cfLXMlKig2EAlBgACfgyRuHghTMq2YoF4wSkpJJUBnjejmAknYtGy2fU2H7rsTOK
BNT0wMUQ10C1eJnYoR7QgyWw3zZOAdo1ivBKagMvdyl+Bz99FrvGQVj4F0EYrvam
7Ee+lNesGfJm6CHuN9eC2uxWNJCxCgr0C+lhnIW6c/23CtbUTqC79ZbBpSqxLUth
AsBLPhYADl31nkFiJKO/r7Ewdwhwahl3lkCxnigBAeUnm0VICtmlNBdu
-----END CERTIFICATE-----`)
	_ = caCertBytes
	clientCert := []byte(`-----BEGIN CERTIFICATE-----
MIIDYTCCAkkCFCDB7NvTBQ9gYdkoDpt8u39faXdbMA0GCSqGSIb3DQEBCwUAMGsx
CzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1TYW4g
RnJhbmNpc2NvMQ4wDAYDVQQKDAVNeU9yZzEPMA0GA1UECwwGTXlEZXB0MQ4wDAYD
VQQDDAVteS1jYTAeFw0yNDA3MjIxMTI0MTRaFw0yNTA3MjIxMTI0MTRaMG8xCzAJ
BgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1TYW4gRnJh
bmNpc2NvMQ4wDAYDVQQKDAVNeU9yZzEPMA0GA1UECwwGTXlEZXB0MRIwEAYDVQQD
DAlteS1jbGllbnQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCVdkCr
/OCEcWLtEAELYscsm7iax5QC7NlZXAT0pLRD1PothuqF+XdId2cm45O1ygzX2joB
OlXH84FhJc/1NBS7MAhfaLvWUcZb8dUbnPXxQwn350YA+hb9gAWoHABeGi5l0wtD
axqXofYfyogeteBIdxib73thWYo3Vma/MOLZA6fNPJqKZM47KzSG70EwHZnCz2rn
RI2EEn8ffodyboKOAQFIfAM8iEj4IO64Y573ko5qBt2M999ulNd98jDJK9pfGrHT
tMNOeQHAg5HJkTLT3KHRMednmixOJEE9UZCPuwBwWAsTBte1ppOShM1Bmt57zIwF
5+saSZvKG30vemkVAgMBAAEwDQYJKoZIhvcNAQELBQADggEBANW9ojqugbe0BfFA
WgHxZhndWE7UfqH6sMvjk59PEn2a6xfuCJKz8PDxvXU3hCKsdf1WQ9I1GpbUYTk7
V6LVZyiXLZNnJYlkGTHgqGvo3mLY3ZCGYoPxeml+bruISG57oEpM7Xk7bD3ecMWd
3I2fEdy5BMNaQNRKQCmsOgVTvbAVWmAv4r4qc2ZVG/gFf5O1kP+Jfl6SWhC0VOoJ
YQv+VD8uDrVh/9vfwtVP39ycLgq2jtejODlbDynTdg+YPshFGvxm++plYO1EJszZ
gUjpJYbIxzCMQefC0DxyUvzN/6B+r+dzzUx9I8KSbIVJQgGUI1avLafVhWk74zCL
3G2VRMs=
-----END CERTIFICATE-----`)
	_ = clientCert
	clientKey := []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCVdkCr/OCEcWLt
EAELYscsm7iax5QC7NlZXAT0pLRD1PothuqF+XdId2cm45O1ygzX2joBOlXH84Fh
Jc/1NBS7MAhfaLvWUcZb8dUbnPXxQwn350YA+hb9gAWoHABeGi5l0wtDaxqXofYf
yogeteBIdxib73thWYo3Vma/MOLZA6fNPJqKZM47KzSG70EwHZnCz2rnRI2EEn8f
fodyboKOAQFIfAM8iEj4IO64Y573ko5qBt2M999ulNd98jDJK9pfGrHTtMNOeQHA
g5HJkTLT3KHRMednmixOJEE9UZCPuwBwWAsTBte1ppOShM1Bmt57zIwF5+saSZvK
G30vemkVAgMBAAECggEAcBOUu2ONGMPN4ua1YcxYfuLms2olW2wwMAoIzUsUwija
0XjyNDS1denTuB2/jfpNVy+Vf4Y2/RFkW2z3XHAJe7SxEpp/AF+h1yCpJWO2KYyT
1QngPKtMwhtWIpGc1PPdBw4SzCNsdXhGD+DX4e+Ql8Z29bfHVWDHfGeV9Ji8Au7N
bAvJjvUYIw7HBPUq7EBkfLJleLg5v3KVSTkHvjtf7K4jOVU+b+pg4qnsEw9Pj9v8
D7vFKrLHFwf2W2uENnNTySHGYUCU0Vm40WkSoeVYyOuc1h+XDCvAHft3r794oroM
Jr4XliY9T0fWDD67YX6S/Cx4d4NteiuIioODeFgwAQKBgQDF1iLO+HpzUlDZ4Jcx
Mvu3glE/Unr36VkDMuBWwOTr+B5QluyLrmBF4NnHY0Nyqa3UwZstDMrw7HkT46a4
gRvU6ptFza+MTxhnn+lcbeJPE52x2nGvKlQGrqMcriXReaX3oVagD3WTKiD1vFYP
1qIB4X+vqRF/YfkFE9cZxNR6lQKBgQDBZ0eGALDEYmDdfCYOEH5FHdOMP1PCvQp4
/MEbhvj9r3EBavaRfbDnuhTcQLAkI/5AKRFGMx8g4cuuajxNeP1XB0jKZeqYdjIq
onJuT25oKBK7x5haWa5EPnhR93E2zXkZI/dM0fzEaX1WYxjnrG5Q23r4GK5jZdZ4
rPPta0EUgQKBgQC3Avq8YBxWpiVpCFyVBMba4dDrNQ/QWqsfTGc/mb2rlKHmh1dX
d/5TZkfQLUFtxw2prVgxeo4aBYeUIJpQQA9RDZ6KGlZ1A45d/g5QlM4vvMO6jYtx
MUT90XvOwkL13wTraPLLqsFnXCeVa55plHHWL5aBF3O6VRWZ3tqzWeP9aQKBgEYh
ARpiHbbYRW+KmPH4oRDG4/Ky89hlW+rLG+qzYo36k+uDsazH+uHL48yJ2FUCiCsT
uSPPXbY6qfSwqPUerh5kkcxycEKgeUhkZ0IAo3Q5M7HLij8YzcwJKu/t3auVjhfD
puTAL/u4lK5CeMFpEQdYzpovuOxp/P79F+Y7QfoBAoGBALqDIODKduaSMnQqHt5q
rq7qic/scGRARaFeSLWwFg7jn46aj3QVmfIn8U+HDqdXrBqJ9oe94a7U1rGSlLOP
+CxeNJHV05ekDiT9tj7GcQ3v/JXJZ9WE/sazbOUNlM4Ugbk+N1lZ2BQ7mwqz5zAE
79NN+xCaA8b959BN8JD6SsFd
-----END PRIVATE KEY-----`)
	_ = clientKey

	url := "amqps://guest:guest@10.1.9.66:5671/"
	// url = "amqp://guest:guest@10.1.9.66:5672/"
	// 1. 构建实例
	rabbitMQ, err := rmq.NewRabbitMQ(url, 5, caCertBytes, clientCert, clientKey)
	// rabbitMQ, err := rmq.NewRabbitMQ(url, 5, nil, nil, nil)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	exchangeName := "exchange001"
	// 2. exchange参数
	exchangeOptions := rmq.ExchangeOptions{
		Type:    "fanout",
		Durable: true,
	}

	// 3. 添加生产者
	rabbitMQ.AddProducer(exchangeName, exchangeOptions)

	// 4. 连接
	err = rabbitMQ.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	i := 0
	for {
		// 5. 发送消息
		log.Printf("Push message: %v\n", i)
		err = rabbitMQ.Publish(exchangeName, "topic1", []byte(fmt.Sprintf("Hello, World: %v", i)))
		if err != nil {
			log.Printf("Failed to publish message: %s\n", err)
			time.Sleep(time.Second * 1)
			continue
		}

		// err = rabbitMQ.Publish("exchange001", "topic2", []byte(fmt.Sprintf("Hello, World: %v", i*10)))
		// if err != nil {
		// 	log.Fatalf("Failed to publish message: %s\n", err)
		// }
		i += 1
		time.Sleep(time.Second * 1)

	}
}
