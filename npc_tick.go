package rogue

import (
	"fmt"
	"math/rand"

	log "github.com/Sirupsen/logrus"
)

// Tick ...
func (n *Npc) Tick() {
	n.Age++

	n.Hunger++
	n.Tiredness++
	n.Thirst++

	fmt.Println("[tick]", n.Name, n.Age)

	if n.isTired() && !n.hasPlanned("sleep") {
		log.Printf("%s is feeling tired. tiredness = %d, cap = %d", n.Name, n.Tiredness, n.tirednessCap())
		n.planAction("sleep")
	}

	if n.isHungry() {

		// auto eat some food in inventory instead of looking for food, if possible
		itemIdx, err := n.tryFindItemTypeInInventory("food")
		if err == nil {
			item := n.removeFromInventory(itemIdx)

			prevHunger := n.Hunger

			// eat item: reduce hunger by some amount from the food eaten
			n.Hunger -= item.Energy
			if n.Hunger < 0 {
				n.Hunger = 0
			}

			energyDiff := prevHunger - n.Hunger
			log.Printf("%s ate %s (-%d hunger)", n.Name, item.Name, energyDiff)
		}

		if n.isHungry() && !n.hasPlanned("find-food") {
			log.Printf("%s is feeling hungry (%d hunger)", n.Name, n.Hunger)
			n.planAction("find-food")
		}
	}

	if n.isThirsty() {

		// auto eat some food in inventory instead of looking for food, if possible
		itemIdx, err := n.tryFindItemTypeInInventory("drink")
		if err == nil {
			item := n.removeFromInventory(itemIdx)

			prevThirst := n.Thirst

			// eat item: reduce hunger by some amount from the food eaten
			n.Thirst -= item.Energy
			if n.Thirst < 0 {
				n.Thirst = 0
			}

			energyDiff := prevThirst - n.Thirst
			log.Printf("%s drank %s (-%d thirst)", n.Name, item.Name, energyDiff)
		}
		if n.isThirsty() && !n.hasPlanned("find-water") {
			log.Printf("%s is feeling thirsty (%d thirst)", n.Name, n.Thirst)
			n.planAction("find-water")
		}
	}

	if !n.isTired() && !n.isHungry() && !n.isThirsty() {
		// when basic needs is resolved, randomly pick something
		// that would help improve situation for the npc

		if !n.hasPlanned("dig-hole") && len(island.withinRadius("rabbit hole", 30, n.Position)) == 0 {
			log.Printf("%s decided to dig a hole (shelter)", n.Name)
			n.planAction("dig-hole")
		}
	}

	// select one action to be doing next
	if len(n.CurrentAction) == 0 && len(n.PlannedActions) > 0 {
		// shuffle actions
		if len(n.PlannedActions) > 1 {
			shuffleActionSlice(n.PlannedActions)
		}

		// pick first
		n.CurrentAction = n.PlannedActions[0]
		n.PlannedActions = n.PlannedActions[1:]
		n.TimeSpentOnCurrentAction = 0

		log.Println(n.Name, "started to", n.CurrentAction)
	}

	n.performCurrentAction()
}

func (n *Npc) performCurrentAction() {
	if len(n.CurrentAction) == 0 {
		return
	}

	status := false
	switch n.CurrentAction {
	case "find-food":
		status = n.performFindFood()
	case "find-water":
		status = n.performFindWater()
	case "sleep":
		status = n.performSleep()
	case "dig-hole":
		status = n.performDigHole()
	default:
		panic(fmt.Errorf("Cant perform unknown action: %s", n.CurrentAction))
	}

	if status == true {
		log.Println(n.Name, "finished performing", n.CurrentAction)
		n.CurrentAction = ""
	}
}

// shuffle slice, without allocations
func shuffleActionSlice(p []string) {

	for i := range p {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}
