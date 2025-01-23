package alsa

func ExampleOpen() {
	card, err := Open(0)
	if err != nil {
		panic("error opening card 0: " + err.Error())
	}
	println("opened card ", card)
}

func ExampleOpenDriver() {
	card, err := OpenDriver("snd_bcm2835")
	if err != nil {
		panic("error opening snd_bcm2835: " + err.Error())
	}
	println("opened card ", card)
}
