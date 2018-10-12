package main

import "fmt"
import "unsafe"
import "math"
type F1 struct {
    Age int
    Name string
}



func main() {
    x := [15]int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
    fmt.Printf("%d\n", x[1:2])

    y := make([]int, 15)
    //copy(y[0:], x[0:])
    fmt.Println(y)


    //copy(y[0:], x[5:])
    //fmt.Println(y)

    fmt.Println("----------")
    fmt.Println(x[:5])
    fmt.Println(x[5:])
    fmt.Println("----------")

    var st F1
    st.Age = 12
    st.Name ="Loooooooooooooooooooooooooooooooooong name"
    fmt.Println(unsafe.Sizeof(st))

    fmt.Println(Round(1.2345, 0.05))

    j := make([]int, 5, 5)
    ji := []int{1}
    copy(j,ji)
    fmt.Println("----------")
    j = append(j, 2)
    j = append(j, 3)
    fmt.Println(j)
    fmt.Println(len(j))
}

func Round(x, unit float64) float64 {
    return math.Round(x/unit) * unit
}
