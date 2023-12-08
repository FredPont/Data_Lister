// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Written by Frederic PONT.
//(c) Frederic Pont 2023

package pogrebdb

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"
)

// ByteSliceToStringslice convert []byte to []float64
func ByteSliceToStringslice(byteArray []byte) []float64 {

	// Convert []byte back to []float64
	FloatArray := make([]float64, len(byteArray)/8)
	for i := 0; i < len(byteArray); i += 8 {
		bits := binary.LittleEndian.Uint64(byteArray[i : i+8])
		FloatArray[i/8] = math.Float64frombits(bits)
	}
	return FloatArray
}

// FloatSliceToStringslice convert []float64 to []string
func FloatSliceToStringslice(FloatArray []float64) []string {

	// Convert []float64 to []string
	stringArray := make([]string, len(FloatArray))
	for i, f := range FloatArray {
		stringArray[i] = strconv.FormatFloat(f, 'f', -1, 64)
	}
	return stringArray
}

func ByteSliceToRow(byteArray []byte) string {
	// Convert []byte back to []float64
	FloatArray := ByteSliceToStringslice(byteArray)
	// Convert []float64 to []string
	stringArray := FloatSliceToStringslice(FloatArray)
	// Convert []string to string
	row := ""
	for _, s := range stringArray {
		row += s + "	"
	}
	return strings.TrimSpace(row) // remove last tabulation
}

// ByteToString convert []byte to a string
func ByteToString(bt []byte) string {
	return string(bt)

}

// StringToByte  convert a string to a []byte
func StringToByte(str string) []byte {
	return []byte(str)
}

// IntToBytes convert int to []byte
func IntToBytes(i int64) []byte {
	//b := make([]byte, 4)
	//binary.LittleEndian.PutUint32(b, uint32(i))
	buf := make([]byte, binary.MaxVarintLen64) // make a byte slice with enough capacity
	n := binary.PutVarint(buf, i)              // encode the int64 value into the byte slice
	b := buf[:n]                               // slice the byte slice to the actual length

	return b
}

func ByteToInt(b []byte) int64 {
	// Convert the byte slice to an uint64 value using binary.BigEndian.Uint64
	//return int64(binary.BigEndian.Uint64(b))
	// var ret uint64
	// buf := bytes.NewBuffer(b)
	// binary.Read(buf, binary.BigEndian, &ret)
	// return int64(ret)

	x, _ := binary.Varint(b) // decode the byte slice to an int64 value
	return x
}

func IntToString(i int64) string {
	return strconv.FormatInt(i, 10) // base 10 for decimal
}
