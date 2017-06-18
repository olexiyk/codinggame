<?php

/**
 * Idea
 */
class Game
{

    const ZOMBIE_SPEED = 400;
    const GUNMAN_SPEED = 1000;
    const GUNMAN_AREA = 2000;

    /**
     * @param Point $gunman
     * @param Points $humans
     * @param Points $zombies
     * @return Point
     */
    function nextTarget($gunman, $humans, $zombies)
    {
        return $this->closestZombieToClosestPossibleToSaveHumanPursue($gunman, $humans, $zombies);
    }

    public static function log($message)
    {
        error_log(var_export($message, true));
    }

    /**
     * @param Point $human
     * @param Point $zombie
     * @param Point $gunman
     * @return bool
     */
    public function possibleToSave(Point $human, Point $zombie, Point $gunman)
    {
        $stepsBetweenHumanAndZombie = $this->stepsBetween(
            $human,
            $zombie,
            self::ZOMBIE_SPEED
        );
        $stepsBetweenHumanAndGunman = $this->stepsBetweenGunmanAnd($human, $gunman);

        if ($stepsBetweenHumanAndZombie < $stepsBetweenHumanAndGunman) {
            return false;
        }
        return true;
    }

    public function stepsBetween(Point $pointFrom, Point $pointTo, $speed)
    {
        return $pointFrom->distanceTo($pointTo) / $speed;
    }

    /**
     * 57220 points
     *
     * @param Point $gunman
     * @param Points $humans
     * @param Points $zombies
     * @return Point
     */
    public function closestZombieToClosestPossibleToSaveHuman($gunman, $humans, $zombies)
    {
        $closestHuman = $humans->closestTo($gunman);
        while (!$this->possibleToSave(
                $closestHuman,
                $zombies->closestTo($closestHuman),
                $gunman
            ) && count($humans) > 1) {
            $this->log("Removing closest human " . $closestHuman);

            $humans->removeClosestTo($gunman);
            $closestHuman = $humans->closestTo($gunman);
        }

        $this->log("Now closest human is " . $closestHuman);

        $closestZombie = $zombies->closestTo($closestHuman);
        return $closestZombie;
    }

    /**
     * 37930 points
     *
     * @param Point $gunman
     * @param Points $humans
     * @param Points $zombies
     * @return Point
     */
    public function closestZombieToClosestPossibleToSaveHumanAimBetween($gunman, $humans, $zombies)
    {
        $closestHuman = $humans->closestTo($gunman);
        while (!$this->possibleToSave(
                $closestHuman,
                $zombies->closestTo($closestHuman),
                $gunman
            ) && count($humans) > 1) {
            $this->log("Removing closest human " . $closestHuman);

            $humans->removeClosestTo($gunman);
            $closestHuman = $humans->closestTo($gunman);
        }

        $this->log("Now closest human is " . $closestHuman);

        $closestZombie = $zombies->closestTo($closestHuman);
        return $closestZombie->pointBetween($gunman);
    }

    /**
     * 40120 points
     *
     * @param Point $gunman
     * @param Points $humans
     * @param Points $zombies
     * @return Point
     */
    public function closestZombieToFarthestPossibleToSaveHuman($gunman, $humans, $zombies)
    {
        $human = $humans->farthestTo($gunman);
        while (!$this->possibleToSave(
                $human,
                $zombies->closestTo($human),
                $gunman
            ) && count($humans) > 1) {
            $this->log("Removing farthest human " . $human);

            $humans->removeFarthestTo($gunman);
            $human = $humans->farthestTo($gunman);
        }

        $this->log("Now farthest human to save is " . $human);

        $closestZombie = $zombies->closestTo($human);
        return $closestZombie;
    }

    /**
     * 53400 points
     *
     * @param Point $gunman
     * @param Points $humans
     * @param Points $zombies
     * @return Point
     */
    public function closestZombieToClosestPossibleToSaveHumanPursue($gunman, $humans, $zombies)
    {
        $closestHuman = $humans->closestTo($gunman);
        while (!$this->possibleToSave(
                $closestHuman,
                $zombies->closestTo($closestHuman),
                $gunman
            ) && count($humans) > 1) {
            $this->log("Removing closest human " . $closestHuman);

            $humans->removeClosestTo($gunman);
            $closestHuman = $humans->closestTo($gunman);
        }

        $this->log("Now closest human is " . $closestHuman);

        $closestZombie = $zombies->closestTo($closestHuman);
        $stepsToHuman = $this->stepsBetweenGunmanAnd($closestHuman, $gunman);
        return $closestZombie->futurePosition($closestHuman, $stepsToHuman, self::ZOMBIE_SPEED);
    }

    /**
     * @param Point $human
     * @param Point $gunman
     * @return float|int
     */
    public function stepsBetweenGunmanAnd(Point $human, Point $gunman)
    {
        $stepsBetweenHumanAndGunman = $this->stepsBetween(
            $human,
            $gunman,
            self::GUNMAN_SPEED
        );
        $stepsBetweenHumanAndGunman -= self::GUNMAN_AREA / self::GUNMAN_SPEED;
        return $stepsBetweenHumanAndGunman < 0 ? 0 : $stepsBetweenHumanAndGunman;
    }
}

class Point
{
    private $x;
    private $y;
    private $id;

    /**
     * Point constructor.
     * @param int $x
     * @param int $y
     * @param int $id
     */
    public function __construct($x, $y, $id = 0)
    {
        $this->x = (int)$x;
        $this->y = (int)$y;
        $this->id = $id;
    }

    /**
     * @param Point $point
     * @return float
     */
    public function distanceTo(Point $point)
    {
        return sqrt(
            pow($this->x - $point->x, 2) +
            pow($this->y - $point->y, 2)
        );
    }

    public function pointBetween(Point $point)
    {
        return new Point(
            ($this->x + $point->x) / 2,
            ($this->y + $point->y) / 2
        );
    }

    public function futurePosition(Point $toPoint, $steps, $speed)
    {
        $xDiff = ($toPoint->x - $this->x) / $this->distanceTo($toPoint) * $steps * $speed;
        $yDiff = ($toPoint->y - $this->y) / $this->distanceTo($toPoint) * $steps * $speed;
        Game::log([
            'steps' => $steps,
            'xDiff' => $xDiff,
            'yDiff' => $yDiff,
        ]);

        return new Point(
            $this->x + $xDiff,
            $this->y + $yDiff
        );
    }

    /**
     * @return int
     */
    public function getId()
    {
        return $this->id;
    }

    public function __toString()
    {
        return $this->x . ' ' . $this->y . "\n";
    }
}

class Points implements Countable
{

    /**
     * @var Point[]
     */
    private $points = array();
    private $ordered = false;

    public function addPoint(Point $point)
    {
        $this->points[] = $point;
        $this->ordered = false;
        return $this;
    }

    private function orderRelatively(Point $toPoint)
    {
//        if (!$this->ordered) {
        usort($this->points, function ($p1, $p2) use ($toPoint) {
            return $toPoint->distanceTo($p1) - $toPoint->distanceTo($p2);
        });
        $this->ordered = true;
//        }
    }

    public function nOrdered(Point $toPoint, $index)
    {
        $this->orderRelatively($toPoint);
        return $this->points[$index];
    }

    public function closestTo(Point $toPoint)
    {
        return $this->nOrdered($toPoint, 0);
    }

    public function farthestTo(Point $toPoint)
    {
        return $this->nOrdered($toPoint, count($this->points) - 1);
    }

    public function removeNOrdered(Point $toPoint, $index)
    {
        $this->orderRelatively($toPoint);
        array_splice($this->points, $index, 1);
    }

    public function removeClosestTo($toPoint)
    {
        $this->removeNOrdered($toPoint, 0);
    }

    public function removeFarthestTo($toPoint)
    {
        $this->removeNOrdered($toPoint, count($this->points) - 1);
    }

    public function count()
    {
        return count($this->points);
    }
}

while (TRUE) {
    fscanf(STDIN, "%d %d",
        $x,
        $y
    );
    $gunman = new Point($x, $y);
    $humans = new Points();
    $zombies = new Points();
    fscanf(STDIN, "%d",
        $humanCount
    );
    for ($i = 0;
         $i < $humanCount;
         $i++) {
        fscanf(STDIN, "%d %d %d",
            $humanId,
            $humanX,
            $humanY
        );
        $humans->addPoint(new Point($humanX, $humanY, $humanId));
    }

    fscanf(STDIN, "%d",
        $zombieCount
    );
    for ($i = 0; $i < $zombieCount; $i++) {
        fscanf(STDIN, "%d %d %d %d %d",
            $zombieId,
            $zombieX,
            $zombieY,
            $zombieXNext,
            $zombieYNext
        );
        $zombies->addPoint(new Point($zombieX, $zombieY, $zombieId));
    }

    $game = new Game();
    $closestZombie = $game->nextTarget($gunman, $humans, $zombies);

    echo($closestZombie); // destination coordinates
}
