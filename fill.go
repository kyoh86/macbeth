package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type fillFlagsStruct struct{}

var (
	fillFlags   fillFlagsStruct
	fillCommand = &cobra.Command{
		Use:   "fill",
		Short: "Fills all or parts of a region with a specific block",
		Args:  cobra.RangeArgs(7, 10),
		RunE: func(_ *cobra.Command, args []string) error {
			params := splitCoordinates(args)
			var remainParams [][]string
			if len(params) > 9 && params[8] == "replace" {
				for _, repl := range strings.Split(params[9], ",") {
					remainParams = append(remainParams, []string{params[6], params[7], params[8], repl})
				}
				// NOTE: unsupported replaceDataValue
			} else {
				remainParams = [][]string{params[7:]}
			}
			fromX, err := parsePoint(params[0])
			if err != nil {
				return fmt.Errorf("parse from-x: %s", err)
			}
			fromY, err := parsePoint(params[1])
			if err != nil {
				return fmt.Errorf("parse from-y: %s", err)
			}
			fromZ, err := parsePoint(params[2])
			if err != nil {
				return fmt.Errorf("parse from-z: %s", err)
			}
			toX, err := parsePoint(params[3])
			if err != nil {
				return fmt.Errorf("parse to-x: %s", err)
			}
			toY, err := parsePoint(params[4])
			if err != nil {
				return fmt.Errorf("parse to-y: %s", err)
			}
			toZ, err := parsePoint(params[5])
			if err != nil {
				return fmt.Errorf("parse to-z: %s", err)
			}
			if fromX.Relative != toX.Relative {
				return fmt.Errorf("x-axis relativity is not equal: (from: %v, to: %v)", fromX.Relative, toX.Relative)
			}
			if fromY.Relative != toY.Relative {
				return fmt.Errorf("y-axis relativity is not equal: (from: %v, to: %v)", fromY.Relative, toY.Relative)
			}
			if fromZ.Relative != toZ.Relative {
				return fmt.Errorf("z-axis relativity is not equal: (from: %v, to: %v)", fromZ.Relative, toZ.Relative)
			}
			xChunks := blockChunks(chunk32Blocks(countBlocks(fromX.Value, toX.Value)))
			yChunks := blockChunks(chunk32Blocks(countBlocks(fromY.Value, toY.Value)))
			zChunks := blockChunks(chunk32Blocks(countBlocks(fromZ.Value, toZ.Value)))
			x := fromX.Value
			for _, xSize := range xChunks {
				y := fromY.Value
				for _, ySize := range yChunks {
					z := fromZ.Value
					for _, zSize := range zChunks {
						fx := Point{Value: x, Relative: fromX.Relative}
						fz := Point{Value: z, Relative: fromZ.Relative}
						tx := Point{Value: x + xSize - 1, Relative: fromX.Relative}
						tz := Point{Value: z + zSize - 1, Relative: fromZ.Relative}
						fmt.Printf("tickingarea add %s 0 %s %s 0 %s temporary-filling\n",
							fx,
							fz,
							tx,
							tz,
						)
						for _, remParam := range remainParams {
							fmt.Printf("fill %s %s %s %s %s %s %s\n",
								fx,
								Point{Value: y, Relative: fromY.Relative},
								fz,
								tx,
								Point{Value: y + ySize - 1, Relative: fromY.Relative},
								tz,
								strings.Join(remParam, " "),
							)
						}
						z = z + zSize
					}
					y = y + ySize
				}
				x = x + xSize
			}
			return nil
		},
	}
)

type Point struct {
	Value    int
	Relative bool
}

func parsePoint(s string) (p Point, _ error) {
	if strings.HasPrefix(s, "~") {
		p.Relative = true
		s = strings.TrimPrefix(s, "~")
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return p, err
	}
	p.Value = v
	return p, nil
}

func (p Point) String() string {
	if p.Relative {
		if p.Value == 0 {
			return "~"
		}
		return fmt.Sprintf("~%d", p.Value)
	}
	return strconv.Itoa(p.Value)
}

func blockChunks(size, count, rem int) (list []int) {
	for i := 0; i < count; i++ {
		list = append(list, size)
	}
	if rem > 0 {
		list = append(list, rem)
	}
	return
}

// 32ブロックを最大とするブロックチャンクのサイズとチャンク数、不足数を求める
func chunk32Blocks(blocks int) (size int, count int, rem int) {
	if blocks == 0 {
		return
	}
	count = blocks / 32
	if blocks%32 > 0 {
		count++
	}
	size = int(math.Ceil(float64(blocks) / float64(count)))
	if surplus := (size * count) - blocks; surplus > 0 {
		rem = size - surplus
		count--
	}
	return
}

// 指定の座標値c1からc2までのブロック数。両者が同じ場合は1になる値である点に注意
func countBlocks(c1, c2 int) int {
	if c1 > c2 {
		return c1 - c2 + 1
	}
	return c2 - c1 + 1
}

func splitCoordinates(args []string) (ret []string) {
	for _, v := range args {
		at := len(ret)
		for r := strings.LastIndex(v, "~"); r > 0; r = strings.LastIndex(v, "~") {
			ret = append(append(append([]string{}, ret[:at]...), v[r:]), ret[at:]...)
			v = v[:r]
		}
		ret = append(append(append([]string{}, ret[:at]...), v), ret[at:]...)
	}
	return
}

func init() {
	facadeCommand.AddCommand(fillCommand)
}
