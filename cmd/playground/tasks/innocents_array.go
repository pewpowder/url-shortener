package tasks

import "fmt"

/*
Условие задачи:
На входе есть два отсортированных по возрастанию массива чисел:
- suspects (подозреваемые)
- innocents (невиновные)

Необходимо из массива подозреваемых исключить всех невиновных и вернуть результирующий массив.

Примеры:
Input: suspects = [1, 2, 3, 4, 5], innocents = [2, 4]
Output: [1, 3, 5]

Input: suspects = [3, 80, 123, 421, 936], innocents = [80, 936]
Output: [3, 123, 421]

Оба массива отсортированы по возрастанию.
*/

func RunInnocentsArray() {
	suspects := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	innocents := []int{2, 3, 5}
	filtered := filter(suspects, innocents)
	fmt.Println(filtered)
}

func filter(suspects, innocents []int) []int {
	ptr1 := 0
	ptr2 := 0
	res := []int{}

	for ptr1 < len(suspects) && ptr2 < len(innocents) {
		if suspects[ptr1] == innocents[ptr2] {
			ptr1++
			ptr2++
			continue
		}

		res = append(res, suspects[ptr1])
		if suspects[ptr1] > innocents[ptr2] {
			ptr2++
		} else {
			ptr1++
		}
	}

	for ptr1 < len(suspects) {
		res = append(res, suspects[ptr1])
		ptr1++
	}

	return res
}
