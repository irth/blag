package main

import (
    "flag"
)

type Config struct {
    Input *string
    Theme *string
    Output *string
}

type BlagPost struct {
    Title string
    Timestamp int
    Author string
    Content string
    Slug string
}


var config Config
func main() {
    config.Input = flag.String("input", "./input/", "Directory where blog posts are stored (in markdown format)");
    config.Output = flag.String("output", "./output/", "Directory where generated html should be stored");
    config.Theme = flag.String("theme", "./theme/", "Directory containing theme files (templates)");
    flag.Parse();
}
