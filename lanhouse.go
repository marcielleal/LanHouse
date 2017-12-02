/*
 * @file lanhouse.go
 * @author Marciel Leal
 * @repository github.com/marcielleal/LanHouse
 * @license GNU license 3.0
 */

/* Sobre comentarios no codigo:
 * Go defende a simplicidade na hora de nomear as coisas, principalmente as
 * variáveis, preferindo nomes curtos a nomes longos, porém, para tornar o 
 * código mais legível, decidi optar pelas boas práticas de programação, 
 * deixando claro o que cada variável significa e o que cada (go)rotina faz,
 * diminuindo a necessidade de comentários.
 */
package main

import(
	"fmt"
	"math/rand"
	"time"
	"sync"
)

const numberOfTeenagers=26
const numberOfComputers=8

type Teenager struct{
	name string
	onlineTime int
	isOnline bool
}

func (t* Teenager) printAccessEvent(){
	fmt.Printf("Adolescente %s. está on-line\n",t.name)
}

func (t* Teenager) printEgressEvent(){
	fmt.Printf("Adolescente %s. liberou a máquina após passar %d minutos\n",t.name,t.onlineTime)
}

func (t* Teenager) printWaitStatus(){
	fmt.Printf("Adolescente %s. está aguardando\n",t.name)
}

func createsTeenagers() ([]Teenager){
	teenagersList:= make([]Teenager, numberOfTeenagers)

	seedToRandom := rand.New(rand.NewSource(time.Now().Unix()))

	listToGenerateRandomNamesOrder:= seedToRandom.Perm(numberOfTeenagers)

	const convertToASCII=65
	const lowerOnlineTimePossible=15
	const sizeOfRangeOfPossiblesOnlineTimes=105

	for i:=0; i<numberOfTeenagers; i++{
		teenagersList[i]=Teenager{
							string(listToGenerateRandomNamesOrder[i] + convertToASCII), // Creates strings in range A..Z
							seedToRandom.Intn(sizeOfRangeOfPossiblesOnlineTimes) + lowerOnlineTimePossible, // Generates numbers between 15 and 120
							i < numberOfComputers } // Only the numberOfComputers-th teenegars will use computers initially
	} 
	return teenagersList 
}

func startsTenagersQueueOnLanhouse()(chan Teenager){
	teenagersList:=createsTeenagers()

	teenagersQueue:=make(chan Teenager, numberOfTeenagers)
	for i:=0; i<numberOfTeenagers; i++{
		teenagersQueue <- teenagersList[i]
		if teenagersList[i].isOnline{
			teenagersList[i].printAccessEvent()
		}else{
			teenagersList[i].printWaitStatus()
		}
	}
	close(teenagersQueue)
	return teenagersQueue
}


func managesComputerUse(teenagersQueue chan Teenager, closingManager * sync.WaitGroup) {
	for{
		teenager, queueIsNotEmpty := <-teenagersQueue
		if queueIsNotEmpty{
			if !teenager.isOnline{
				teenager.printAccessEvent()
			}
			time.Sleep(time.Duration(teenager.onlineTime/2.0)*time.Second)
			teenager.printEgressEvent()
		}else{
			closingManager.Done()
			return
		}
	}
}

func managesLanHouse(closingManager * sync.WaitGroup){
	teenagersQueue:=startsTenagersQueueOnLanhouse()
	for i := 0; i < numberOfComputers; i++ {
		go managesComputerUse(teenagersQueue, closingManager)
	}
	closingManager.Wait()

	fmt.Printf("\nA lan-house está finalmente vazia e todos foram atendidos\n")
}

func main(){
	var closingManager sync.WaitGroup		//Manages the lanhouse closing
	closingManager.Add(numberOfComputers)	

	go managesLanHouse(&closingManager)		
	closingManager.Wait()					//Do not end main until lanhouse is closed
}
