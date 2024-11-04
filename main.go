package main

import (
	"fmt"
	"sync"
	"time"
)

const chunkSize = 20

func newSprites(sb *SpriteBucket, num int) []*Sprite {
	sprites := []*Sprite{}
	for i := 0; i < num; i++ {
		sprites = append(sprites, newSprite(sb))
	}
	return sprites
}

func newSprite(sb *SpriteBucket) *Sprite {
	return &Sprite{
		x:            0,
		y:            0,
		spriteBucket: sb,
	}
}

type Sprite struct {
	x, y         int
	spriteBucket *SpriteBucket
}

// func (s []*Sprite) A() {
// 	//logic here
// 	time.Sleep(time.Duration(s.durationToUpdate))
// 	fmt.Print(s.spriteBucket.name)
// }

func (s *Sprite) Update() {
	//logic here
	time.Sleep(time.Duration(s.spriteBucket.duration))
	fmt.Print(s.spriteBucket.name)
}

type SpriteBucket struct {
	sprites  []*Sprite
	name     string
	wg       sync.WaitGroup
	duration time.Duration
}

// Update a chunk of sprites within this SpriteBucket
func (sb *SpriteBucket) startSpriteRoutine(sprites []*Sprite) {
	defer sb.wg.Done()
	var wg sync.WaitGroup
	wg.Add(len(sprites))
	for _, sprite := range sprites {
		go func() {
			defer wg.Done()
			sprite.Update()
		}()
	}
	wg.Wait()
}

func (sb *SpriteBucket) Update() {

	// Split sprites into chunks and run each chunk in its own goroutine
	for i := 0; i < len(sb.sprites); i += chunkSize {
		end := i + chunkSize
		if end > len(sb.sprites) {
			end = len(sb.sprites) // Handle last chunk which may be smaller
		}
		sb.wg.Add(1)
		go sb.startSpriteRoutine(sb.sprites[i:end])
	}

	sb.wg.Wait() // Wait for all chunk goroutines in this bucket to complete
}

type RenderManager struct {
	buckets map[int]*SpriteBucket
	wg      sync.WaitGroup
}

func newRenderManager() *RenderManager {
	return &RenderManager{
		buckets: make(map[int]*SpriteBucket),
	}
}

func (rm *RenderManager) AddBucketOnTop(bucket *SpriteBucket) *RenderManager {
	num := len(rm.buckets)
	rm.buckets[num] = bucket
	return rm
}

func (rm *RenderManager) UpdateAll() {
	for i := range len(rm.buckets) {
		if bucket, ok := rm.buckets[i]; ok {
			rm.wg.Add(1)
			go func() {
				defer rm.wg.Done()
				bucket.Update()
			}()
		}
	}

	rm.wg.Wait() // wait for all buckets to finish updating
}

func main() {
	fmt.Println("Loading assets.")
	var (
		frontBucket = &SpriteBucket{
			name:     "front",
			duration: 0 * time.Millisecond,
		}
		midBucket = &SpriteBucket{
			name:     "mid",
			duration: 0 * time.Millisecond,
		}
		backBucket = &SpriteBucket{
			name:     "back",
			duration: 0 * time.Millisecond,
		}
		RenderManager = newRenderManager().AddBucketOnTop(backBucket).AddBucketOnTop(midBucket).AddBucketOnTop(frontBucket)
	)

	// update takes sprites/buffersize * longest spritebucket duration

	frontBucket.sprites = newSprites(frontBucket, 20)
	midBucket.sprites = newSprites(midBucket, 20)
	backBucket.sprites = newSprites(backBucket, 20)

	fmt.Println("Assets loaded.")

	start := time.Now()
	RenderManager.UpdateAll()
	done := time.Since(start)

	fmt.Println("\nUpdate completed. Took ", done)
}
