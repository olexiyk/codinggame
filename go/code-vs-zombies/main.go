/*
The Goal

Destroy zombies quickly to earn points and make sure to keep the humans alive to get the highest score that you can manage.
Rules

The game is played in a zone 16000 units wide by 9000 units high. You control a man named Ash, wielding a gun that lets him kill any zombie within a certain range around him.

Ash works as follows:
Ash can be told to move to any point within the game zone by outputting a coordinate X Y. The top-left point is 0 0.
Each turn, Ash will move exactly 1000 units towards the target coordinate, or onto the target coordinates if he is less than 1000 units away.
If at the end of a turn, a zombie is within 2000 units of Ash, he will shoot that zombie and destroy it. More details on combat further down.

Other humans will be present in the game zone, but will not move. If zombies kill all of them, you lose the game and score 0 points for the current test case.

Zombies are placed around the game zone at the start of the game, they must be destroyed to earn points.

Zombies work as follows:
Each turn, every zombie will target the closest human, including Ash, and step 400 units towards them. If the zombie is less than 400 units away, the human is killed and the zombie moves onto their coordinate.
Two zombies may occupy the same coordinate.

The order in which actions happens in between two rounds is:
Zombies move towards their targets.
Ash moves towards his target.
Any zombie within a 2000 unit range around Ash is destroyed.
Zombies eat any human they share coordinates with.

Killing zombies earns you points. The number of points you get per zombie is subject to a few factors.

Scoring works as follows:
A zombie is worth the number of humans still alive squared x10, not including Ash.
If several zombies are destroyed during on the same round, the nth zombie killed's worth is multiplied by the (n+2)th number of the Fibonnacci sequence (1, 2, 3, 5, 8, and so on). As a consequence, you should kill the maximum amount of zombies during a same turn.

Note: You may activate gory mode in the settings panel () if you have the guts for it.

Victory Conditions
You destroy every zombie on screen with at least 1 other living human remaining.

Lose Conditions
The zombies kill every human other than you.
Expert Rules

The coordinate system of the game uses whole numbers only. If Ash or a zombie should land in a non whole coordinate, that coordinate is rounded down.

For example, if a zombie were to move from X=0, Y=0 towards X=500, Y=500, since it may only travel 400 units in one turn it should land on X=282.843, Y=282.843 but will in fact land on X=282, Y=282.

To help, each zombie's future coordinates will be sent along side the current coordinates.
Example

The player starts at position X=8043, Y=976. Two zombies are close by but the player decides to go near a human at X=0, Y=4500.

Turn 1
Action: 0 3433.
Zombies 0 and 1 both aim for the player.

Turn 2
Action: 0 3833.
The player gets closer to zombie 1 and further away from zombie 0.

Turn 3
Action: 0 4233.
Zombie 1 enters player's range and is destroyed!




Note

Don’t forget to run the tests by launching them from the “Test cases” window. You do not have to pass all tests to enter the leaderboard. Each test you pass will earn you some points.

Warning: the tests provided are similar to the validation tests used to compute the final score but remain different. This is a “hardcoding” prevention mechanism. Harcoded solutions will not get any points.

Your score is computed from the total points earned across all test cases.

Do not hesitate to switch to debug mode () if you get stuck. In debug mode, hover over a zombie or human to see their coordinates.
Game Input

The program must, within an infinite loop, read the contextual data from the standard input (human and zombie positions) and provide to the standard output the desired instruction.
Input for one game turn
Line 1: two space-separated integers x and y, the coordinate of your character.

Line 2: one integer humanCount, the amount of other humans still alive.

Next humanCount lines : three space-separated integers humanId, humanX & humanY, the unique id and coordinates of a human.

Next line: one integer zombieCount, the amount of zombies left to destroy.

Next zombieCount lines: five space-separated integers zombieId, zombieX, zombieY, zombieXNext & zombieYNext, the unique id, current coordinates and future coordinates of a zombie.

Output for one game turn
A single line: two integers targetX and targetY, the coordinate you want your character to move towards. You may also some text message which will be displayed on screen.
Constraints
0 ≤ x < 16000

0 ≤ y < 9000

1 ≤ humanCount < 100

1 ≤ zombieCount < 100

Response time per game turn ≤ 100ms

*/

package main

import (
	"fmt"
	"math"
	"encoding/json"
	"math/rand"
	"time"
	"sort"
)

const maxX = 16000
const maxY = 9000
const maxStepSize = 1000
const globalTimeLimitInMilliseconds = 100
const maxSteps = 2
const populationSize = 100 * 2

func main() {
	timeChan := time.NewTimer(time.Millisecond * globalTimeLimitInMilliseconds).C

	g := Game{}
	s1 := rand.NewSource(time.Now().UnixNano())
	g.random = rand.New(s1)

	bestPath := new(Path)
Loop:
	for {
		select {
		case <-timeChan:
			fmt.Println("Timer expired")
			break Loop
		default:
			bestPath = g.GetBestPath()
		}
	}
	Log(g.generation)
	Log("Best path")
	Log(bestPath)
}

type Game struct {
	population []*Path
	random     *rand.Rand
	generation uint64
}

func (game *Game) GetBestPath() *Path {
	if len(game.population) < 1 {
		for i := 0; i < populationSize; i++ {
			game.population = append(game.population, game.GenerateRandomPath(maxSteps))
		}
	}

	sort.Sort(ByScore(game.population))
	game.eliminateBadIndividuals()
	game.addNewRandomIndividuals()
	game.crossBreedGoodIndividuals()
	game.evaluatePopulation()
	game.generation++
	return game.doGetBestPath()
}
func (game Game) doGetBestPath() *Path {
	return game.population[0]
}
func (game Game) crossBreedGoodIndividuals() {

}
func (game Game) addNewRandomIndividuals() {

}
func (game Game) eliminateBadIndividuals() {

}

func (game Game) evaluatePopulation() {
	for _, path := range game.population {
		if path.Score < 1 {
			path.Score = rand.Intn(1000)
		}
	}
}

type Path struct {
	Points []*Point `json:"Points"`
	Score  int `json:"Score"`
}

// ByAge implements sort.Interface
type ByScore []*Path

func (population ByScore) Len() int           { return len(population) }
func (population ByScore) Swap(i, j int)      { population[i], population[j] = population[j], population[i] }
func (population ByScore) Less(i, j int) bool { return population[i].Score > population[j].Score }

func (p *Path) addPoint(point *Point) *Path {
	p.Points = append(p.Points, point)
	return p
}

type Point struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}

func (p Point) distanceTo(to Point) float64 {
	return math.Sqrt(math.Pow(p.X-to.X, 2) +
		math.Pow(p.Y-to.Y, 2))
}

func NewPoint(x float64, y float64) *Point {
	p := new(Point)
	p.X = math.Floor(x)
	p.Y = math.Floor(y)
	// points outside of the area are not allowed and truncated
	if p.X > maxX {
		p.X = maxX
	}
	if p.Y > maxY {
		p.Y = maxY
	}
	if p.X < 0 {
		p.X = 0
	}
	if p.Y < 0 {
		p.Y = 0
	}
	return p
}

func Log(v interface{}) {
	data, _ := json.Marshal(v)
	fmt.Printf(string(data) + "\n")
}

func (game Game) GenerateRandomPath(numberOfSteps int) *Path {
	path := new(Path)
	// start anywhere
	startPoint := NewPoint(
		game.random.Float64()*maxX,
		game.random.Float64()*maxY,
	)
	path.addPoint(startPoint)

	for i := 1; i < numberOfSteps; i++ {
		prevPoint := path.Points[len(path.Points)-1]
		path.addPoint(NewPoint(
			prevPoint.X+(game.random.Float64()-0.5)*maxStepSize,
			prevPoint.Y+(game.random.Float64()-0.5)*maxStepSize,
		))
	}
	return path
}

func (game Game) BreedPaths(p1 *Path, p2 *Path) *Path {
	newPath := new(Path)
	for i, point1 := range p1.Points {
		point2 := p2.Points[i]
		average := NewPoint(
			(point1.X+point2.X)/2,
			(point1.Y+point2.Y)/2,
		)
		newPath.addPoint(average)
	}
	return newPath
}
