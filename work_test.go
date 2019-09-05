package future

import (
	"context"
	"log"
	"testing"
	"time"
)

func Test_Wait(t *testing.T) {

	wg := NewGroupWithContext(context.Background())

	wg.Submit(func() (i interface{}, e error) {
		log.Print("abc jajaj ")
		return nil, nil
	})
	wg.Submit(func() (i interface{}, e error) {
		log.Print("abc 2 ")
		return nil, nil
	})

	wg.Wait()

	t.Log("test success ")
}

func Test_Timeout(t *testing.T) {
	wg := NewGroupWithTimeout(1 * time.Second)

	wg.Submit(func() (i interface{}, e error) {
		log.Print("abc jajaj ")
		<-time.NewTimer(2 * time.Second).C
		return nil, nil
	})
	wg.Submit(func() (i interface{}, e error) {
		log.Print("abc 2 ")
		return nil, nil
	})

	if wg.Wait() {
		log.Println("time out")
		t.Log("test success ")
	}

}

func Test_Future(t *testing.T) {

	wg := NewGroup()
	ft1 := wg.Submit(func() (i interface{}, e error) {

		return 1, nil
	})
	ft2 := wg.Submit(func() (i interface{}, e error) {

		return "sad", nil
	})

	v, err := ft1.Get()
	log.Printf("%d %s", v.(int), err)

	v, err = ft2.Get()
	log.Printf("%s %s", v.(string), err)

	t.Log("test success ")
}

func TestWorkgroup_WithSize(t *testing.T) {

	wg := NewGroup().WithSize(3)

	for i := 0; i < 10; i++ {
		index := i
		wg.Submit(func() (i interface{}, e error) {
			<-time.NewTimer(2 * time.Second).C
			log.Printf("%d running ", index)
			return 0, nil
		})
	}

	wg.Wait()

	log.Println("all done")
}

func TestWorkgroup_WithSizeAndTimeout(t *testing.T) {

	wg := NewGroupWithTimeout(5 * time.Second).WithSize(3)

	for i := 0; i < 10; i++ {
		index := i
		wg.Submit(func() (i interface{}, e error) {
			<-time.NewTimer(2 * time.Second).C
			log.Printf("%d running ", index)
			return 0, nil
		})
	}

	if wg.Wait() {
		log.Println("timeout")
		t.Log("success ")
	} else {
		log.Println("all done")
	}

}
