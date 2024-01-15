package main

func main() {
	d := &database{}
	d.build("bench")
	d.clear()
	d.bench()
}
