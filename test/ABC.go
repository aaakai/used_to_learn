package test

func ABC() {
	cha := make(chan struct{})
	chb := make(chan struct{})
	chc := make(chan struct{})
	chClose := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			<-cha
			println("a")
			chb <- struct{}{}
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			<-chb
			println("b")
			chc <- struct{}{}
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			<-chc
			println("c")
			if i == 9 {
				chClose <- struct{}{}
				return
			}
			cha <- struct{}{}
		}
	}()
	cha <- struct{}{}
	<-chClose
	return
}
