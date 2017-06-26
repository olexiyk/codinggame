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
	//"math"
	"encoding/json"
	"math/rand"
	"sort"
	"time"
	"os"
	"math"
)

const maxX = 16000
const maxY = 9000
const maxStepSize = 1000

// it's claimed that timeout is 100ms but 140ms is still tolerated
const globalTimeLimitInMilliseconds = 130
const maxSteps = 20
const populationSize = 100
const randomnessPercent = 0.2
const ashKillingDistance = 2000
const zombieDistancePerStep = 400
const test = true

func main() {
	for {

		g := Game{}

		if !test {
			var x, y int
			fmt.Scan(&x, &y)
			fmt.Fprintf(os.Stderr, "%d %d\n", x, y)

			g.ash = newPoint(x, y)

			var humanCount int
			fmt.Scan(&humanCount)
			fmt.Fprintf(os.Stderr, "%d\n", humanCount)

			g.people = make([]*Point, humanCount)
			for i := 0; i < humanCount; i++ {
				var humanId, humanX, humanY int
				fmt.Scan(&humanId, &humanX, &humanY)
				fmt.Fprintf(os.Stderr, "%d %d %d\n", humanId, humanX, humanY)
				p := newPoint(humanX, humanY)
				g.people[i] = p
			}
			var zombieCount int
			fmt.Scan(&zombieCount)
			fmt.Fprintf(os.Stderr, "%d\n", zombieCount)

			g.zombies = make([]*Point, zombieCount)
			for i := 0; i < zombieCount; i++ {
				var zombieId, zombieX, zombieY, zombieXNext, zombieYNext int
				fmt.Scan(&zombieId, &zombieX, &zombieY, &zombieXNext, &zombieYNext)
				fmt.Fprintf(os.Stderr, "%d %d %d %d %d\n", zombieId, zombieX, zombieY, zombieXNext, zombieYNext)
				p := newPoint(zombieX, zombieY)
				g.zombies[i] = p
			}
		} else {

			g.ash = newPoint(0, 0)
			g.people = append(g.people, newPoint(8250, 4500))
			g.zombies = append(g.zombies, newPoint(8250, 8999))
		}
		s1 := rand.NewSource(time.Now().UnixNano())
		g.random = rand.New(s1)

		bestPath := new(Path)

		timeChan := time.NewTimer(time.Millisecond * globalTimeLimitInMilliseconds).C

	Loop:
		for {
			select {
			case <-timeChan:
				break Loop
			default:
				bestPath = g.GetBestPath()
				Log("Generation number", g.generation)
				//Log("Best path")
				//Log(bestPath)
			}
		}
		Log("Best path", bestPath)
		Log("Best path score ", bestPath.Score)
		fmt.Printf("%d %d\n", bestPath.Points[0].X, bestPath.Points[0].Y) // Your destination coordinates
	}
}

type human Point

type Game struct {
	paths      []*Path
	random     *rand.Rand
	zombies    []*Point
	people     []*Point
	ash        *Point
	generation uint
	run        bool
}

func (game *Game) GetBestPath() *Path {
	// initial paths
	if !game.run {
		game.paths = make([]*Path, populationSize)
		for i := 0; i < populationSize; i++ {
			game.paths[i] = game.generateRandomPath(maxSteps)
		}
		game.run = true
	}
	sort.Sort(byScores(game.paths))
	game.replaceBadIndividualsWithRandom()
	game.crossBreedGoodIndividuals()
	game.evaluatePopulation()
	game.generation++
	return game.doGetBestPath()
}
func (game *Game) doGetBestPath() *Path {
	return game.paths[0]
}
func (game *Game) crossBreedGoodIndividuals() {
	for i := 0; i < populationSize; i += 2 {
		game.paths[i] = game.breedPaths(game.paths[i], game.paths[i+1])
	}
}
func (game *Game) replaceBadIndividualsWithRandom() {
	n := int(populationSize * (1 - randomnessPercent))
	for i := n; i < populationSize; i++ {
		game.paths[i] = game.generateRandomPath(maxSteps)
	}
}
func (game *Game) evaluatePopulation() {
	for _, path := range game.paths {
		game.evaluatePath(path)
	}
}

func (game *Game) evaluatePath(path *Path) {
	zombies := make([]*Point, len(game.zombies))
	copy(zombies, game.zombies)
	people := make([]*Point, len(game.people))
	copy(people, game.people)
	ash := newPoint(game.ash.X, game.ash.Y)

	for _, step := range path.Points {
		if len(zombies) < 1 {
			return
		}

		// 1. Zombies move towards their targets.
		for zi, zombie := range zombies {
			humans := append(people, ash)
			closestHuman := findClosest(zombie, humans)
			zombies[zi] = zombies[zi].moveToDirection(closestHuman, zombieDistancePerStep)
		}
		// 2. Ash moves towards his target.
		ash = step
		// 3. Any from within a 2000 unit range around Ash is destroyed.
		i := 0
		zombiesKilled := 0
		peopleAlive := len(people)
		for _, zombie := range zombies {
			// only left non-killed zombies
			if zombie.distanceTo(ash) > ashKillingDistance {
				zombies[i] = zombie
				i++

				// 4. Zombies eat any human they share coordinates with.
				j := 0
				for _, human := range people {
					if zombie.X != human.X || zombie.Y != human.X {
						people[j] = human
						j++
					} else {
						peopleAlive--
					}
				}
				people = people[:j]
			} else {
				zombiesKilled++
			}
		}
		zombies = zombies[:i]

		// Calculate scores
		if peopleAlive == 0 {
			path.Score = 0
			return
		}
		path.Score += uint64(peopleAlive*peopleAlive*10) * fib[zombiesKilled]
	}
}
func findClosest(from *Point, to []*Point) *Point {
	finder := closestFinder{
		from: from,
		to:   to,
	}
	sort.Sort(byDistance(finder))
	return finder.to[0]
}
func (from Point) moveToDirection(to *Point, maxDistance int) *Point {
	// TODO check if it's < or <=
	if from.distanceTo(to) < maxDistance {
		return *(&to)
	}

	fromMinusTo := newPointUnchecked(
		to.X-from.X,
		to.Y-from.Y,
	)

	norm := fromMinusTo.norm()
	diff := norm / float64(maxDistance)
	point := newPoint(
		int(float64(from.X)+float64(fromMinusTo.X)/diff),
		int(float64(from.Y)+float64(fromMinusTo.Y)/diff),
	)
	return point

	//if from.X-to.X == 0 {
	//	return newPoint(
	//		from.X,
	//		from.Y-maxDistance,
	//	)
	//}
	//
	//angle := float64(from.Y-to.Y) / float64(from.X-to.X)
	//atan := math.Atan(angle)
	//i := float64(maxDistance) * math.Cos(atan)
	//i3 := float64(maxDistance) * math.Sin(atan)
	//return newPoint(
	//	from.X-int(i),
	//	from.Y-int(i3),
	//)

}

func (p *Point) norm() float64 {
	return math.Sqrt(math.Pow(float64(p.X), 2) +
		math.Pow(float64(p.Y), 2))
}

type zombie Point

type closestFinder struct {
	from *Point
	to   []*Point
}

type byDistance closestFinder

func (peopleAndZombie byDistance) Len() int           { return len(peopleAndZombie.to) }
func (peopleAndZombie byDistance) Swap(i, j int)      { peopleAndZombie.to[i], peopleAndZombie.to[j] = peopleAndZombie.to[j], peopleAndZombie.to[i] }
func (peopleAndZombie byDistance) Less(i, j int) bool { return peopleAndZombie.from.distanceTo(peopleAndZombie.to[i]) < peopleAndZombie.from.distanceTo(peopleAndZombie.to[j]) }

type Path struct {
	Points []*Point `json:"Points"`
	Score  uint64 `json:"Score"`
}

// byScores implements sort.Interface
type byScores []*Path

func (population byScores) Len() int           { return len(population) }
func (population byScores) Swap(i, j int)      { population[i], population[j] = population[j], population[i] }
func (population byScores) Less(i, j int) bool { return population[i].Score > population[j].Score }

func (p *Path) addPoint(point *Point) *Path {
	p.Points = append(p.Points, point)
	return p
}

type Point struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

func (p Point) distanceTo(to *Point) int {
	return int(math.Sqrt(math.Pow(float64(p.X-to.X), 2) +
		math.Pow(float64(p.Y-to.Y), 2)))
}

func newPoint(x int, y int) *Point {
	p := newPointUnchecked(x, y)
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

func newPointUnchecked(x int, y int) *Point {
	p := new(Point)
	p.X = x
	p.Y = y
	return p
}

func Log(message string, v interface{}) {
	data, _ := json.Marshal(v)
	fmt.Fprintln(os.Stderr, message+": "+string(data))
}

func (game *Game) generateRandomPath(numberOfSteps int) *Path {
	path := new(Path)
	path.Points = make([]*Point, numberOfSteps)

	prevPoint := game.ash

	d := game.random.Intn(100)
	for i := 0; i < numberOfSteps; i++ {
		nextStep := &Point{}
		if d < 60 {
			closestHuman := findClosest(prevPoint, game.people)
			nextStep = prevPoint.moveToDirection(closestHuman, maxStepSize)
		} else if d < 20 {
			closestZombie := findClosest(prevPoint, game.zombies)
			nextStep = prevPoint.moveToDirection(closestZombie, maxStepSize)
		} else {
			x := game.random.Intn(maxStepSize) - (maxStepSize / 2)
			y := game.random.Intn(maxStepSize) - (maxStepSize / 2)
			nextStep = newPoint(
				prevPoint.X+x,
				prevPoint.Y+y,
			)
		}
		path.Points[i] = nextStep

		if i > 0 {
			prevPoint = path.Points[i-1]
		}

	}
	return path
}

func (game *Game) breedPaths(p1 *Path, p2 *Path) *Path {
	newPath := new(Path)
	newPath.Points = make([]*Point, len(p1.Points))

	for i, point1 := range p1.Points {
		point2 := p2.Points[i]
		average := newPoint(
			(point1.X+point2.X)/2,
			(point1.Y+point2.Y)/2,
		)
		newPath.Points[i] = average
	}
	return newPath
}

// first 100 fibonacci numbers
var fib = [100] uint64{
	0,
	1,
	1,
	2,
	3,
	5,
	8,
	13,
	21,
	34,
	55,
	89,
	144,
	233,
	377,
	610,
	987,
	1597,
	2584,
	4181,
	6765,
	10946,
	17711,
	28657,
	46368,
	75025,
	121393,
	196418,
	317811,
	514229,
	832040,
	1346269,
	2178309,
	3524578,
	5702887,
	9227465,
	14930352,
	24157817,
	39088169,
	63245986,
	102334155,
	165580141,
	267914296,
	433494437,
	701408733,
	1134903170,
	1836311903,
	2971215073,
	4807526976,
	7778742049,
	12586269025,
	20365011074,
	32951280099,
	53316291173,
	86267571272,
	139583862445,
	225851433717,
	365435296162,
	591286729879,
	956722026041,
	1548008755920,
	2504730781961,
	4052739537881,
	6557470319842,
	10610209857723,
	17167680177565,
	27777890035288,
	44945570212853,
	72723460248141,
	117669030460994,
	190392490709135,
	308061521170129,
	498454011879264,
	806515533049393,
	1304969544928657,
	2111485077978050,
	3416454622906707,
	5527939700884757,
	8944394323791464,
	14472334024676221,
	23416728348467685,
	37889062373143906,
	61305790721611591,
	99194853094755497,
	160500643816367088,
	259695496911122585,
	420196140727489673,
	679891637638612258,
	1100087778366101931,
	1779979416004714189,
	2880067194370816120,
	4660046610375530309,
	7540113804746346429,
	12200160415121876738,
	//19740274219868223167,
	//31940434634990099905,
	//51680708854858323072,
	//83621143489848422977,
	//135301852344706746049,
	//218922995834555169026,
}
