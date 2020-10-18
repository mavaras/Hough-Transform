package main

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "math"
    "os"
)

type Image interface {
    Set(int, int, color.Color)
}

func prints(text string) {
    fmt.Println(text)
}

func get_hough_space(img image.Image, height int, width int) draw.Image {
    var (
        y_max = float64(height)
        x_max = float64(width)
        theta_dim = 360
        theta_max = float64(math.Pi)
        r_dim = 700
        r_max = math.Hypot(x_max, y_max)
    )
    hough_space := image.NewGray(image.Rect(0, 0, theta_dim, int(r_max)))
    draw.Draw(
        hough_space,
        hough_space.Bounds(),
        image.NewUniform(color.White),
        image.Point{},
        draw.Src,
    )

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            pixel := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
            if pixel.Y == 255 {
                continue
            }
            x := x - width / 2
            y := y - height / 2
            for itheta := 0; itheta < theta_dim; itheta++ {
                theta := (float64(itheta) * theta_max) / float64(theta_dim)
                _ = theta
                r := float64(height/2) - float64(x)*math.Cos(float64(theta)) + float64(y)*math.Sin(float64(theta))
                ir := (float64(r_dim) * r) / r_max
                _ = ir
                pixel = hough_space.At(int(itheta), int(r)).(color.Gray)
                if pixel.Y > 0 {
                    pixel.Y--
                    hough_space.Set(int(itheta), int(r), pixel)
                }
            }
        }
    }

    return hough_space
}

func get_hough_space_maxs(hough_space draw.Image, height int, width int) (draw.Image, [][]int) {
    var (
        y_max = float64(height)
        x_max = float64(width)
        theta_dim = 360
        r_max = math.Hypot(x_max, y_max)
    )
    hough_space_maxs := image.NewRGBA(image.Rect(0, 0, theta_dim, int(r_max)))
    draw.Draw(
        hough_space_maxs,
        hough_space_maxs.Bounds(),
        image.NewUniform(color.White),
        image.Point{},
        draw.Src,
    )
    var maximums [][]int

    for x := 0; x < theta_dim; x++ {
        for y := 0; y < int(r_max); y++ {
            pixel := hough_space.At(x, y).(color.Gray).Y
            if pixel < 20 {
                maximums = append(maximums, []int{x, y})
                hough_space_maxs.Set(x, y, color.RGBA{255, 0, 0, 255})
            }
        }
    }

    return hough_space_maxs, maximums
}

func convert_hough_to_xy(maximums [][]int, height int, width int) draw.Image {
    res_image := image.NewRGBA(image.Rect(0, 0, width, height))

    for c := 0; c < len(maximums); c++ {
        theta := float64(maximums[c][0])
        rho := float64(maximums[c][1])

        a := math.Cos((theta + 1) * float64(math.Pi / 360.0))
        b := math.Sin((theta + 1) * float64(math.Pi / 360.0))
        x1 := a * rho + 1000 * -b
        y1 := b * rho + 1000 * a
        x2 := a * rho - 1000 * -b
        y2 := b * rho - 1000 * a
        fmt.Printf("%f %f, %f %f\n", x1, y1, x2, y2)
        fmt.Printf("%f %f\n", theta, rho)
        fmt.Printf("%f %f\n\n", a, b)

        for c := 0.0; c > -float64(height)-180.0 && c < float64(height)+180.0; c++ {
            y := math.Cos((theta + 1) * float64(math.Pi / 360.0)) * c - math.Sin(-(theta + 1) * float64(math.Pi / 360.0)) * rho
            x := math.Sin((theta + 1) * float64(math.Pi / 360.0)) * c + math.Cos(-(theta + 1) * float64(math.Pi / 360.0)) * rho
            res_image.Set(int(x), int(y), color.RGBA{205, 105, 55, 255})
        }
    }

    return res_image
}

func save_image(img image.Image, filename string) {
    res_file, err := os.Create(filename)
    if err != nil {
        fmt.Println(err)
        return
    }
    png.Encode(res_file, img)
}

func main() {
    file, err := os.Open("./testimgs/line5.png")

    if err != nil {
        prints("Error: File could not be opened")
        os.Exit(1)
    }

    img, _, err := image.Decode(file)
    file.Close()
    bounds := img.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y

    if err != nil {
        prints("Error: File could not be read")
        os.Exit(1)
    }

    hough_space := get_hough_space(img, height, width)
    hough_space_maxs, maximums := get_hough_space_maxs(hough_space, height, width)
    res_image := convert_hough_to_xy(maximums, height, width)

    save_image(hough_space, "hough.png")
    save_image(hough_space_maxs, "houghmxs.png")
    save_image(res_image, "res.png")

    file.Close()
}
