package main

import (
	"errors"
	"fmt"
)

const BITMASK_7 uint8 = 1
const BITMASK_6 uint8 = 2
const BITMASK_5 uint8 = 4
const BITMASK_4 uint8 = 8
const BITMASK_3 uint8 = 16
const BITMASK_2 uint8 = 32
const BITMASK_1 uint8 = 64
const BITMASK_0 uint8 = 128

const MAX_HEIGHT uint8 = 5
const MAX_COLUMNS uint8 = 6

const PLAYER_ZERO_BOARD = 1
const PLAYER_ONE_BOARD = 2

// DATA STRUCTURE OF "compressed_board"
// col1: uint8 0(b0-5)
// col2: uint8 0(b6-7), 1(b0-3):
// col3: uint8 1(b4-7), 1(b0-1):
// col4: uint8 2(b2-7)
// col5: uint8 3(b0-5)
// col2: uint8 3(b6-7), 1(b0-3):
// col3: uint8 1(b4-7), 1(b0-1):

type comp_board_bitmask struct {
	start_byte uint8
	masks      []uint8
}

var COLUMN_0_BITMASK = comp_board_bitmask{
	start_byte: 0,
	masks:      []uint8{BITMASK_0, BITMASK_1, BITMASK_2, BITMASK_3, BITMASK_4, BITMASK_5},
}

var COLUMN_1_BITMASK = comp_board_bitmask{
	start_byte: 0,
	masks:      []uint8{BITMASK_6, BITMASK_7, BITMASK_0, BITMASK_1, BITMASK_2, BITMASK_3},
}

var COLUMN_2_BITMASK = comp_board_bitmask{
	start_byte: 1,
	masks:      []uint8{BITMASK_4, BITMASK_5, BITMASK_6, BITMASK_7, BITMASK_0, BITMASK_1},
}

var COLUMN_3_BITMASK = comp_board_bitmask{
	start_byte: 2,
	masks:      []uint8{BITMASK_2, BITMASK_3, BITMASK_4, BITMASK_5, BITMASK_6, BITMASK_7},
}

var COLUMN_4_BITMASK = comp_board_bitmask{
	start_byte: 3,
	masks:      []uint8{BITMASK_0, BITMASK_1, BITMASK_2, BITMASK_3, BITMASK_4, BITMASK_5},
}

var COLUMN_5_BITMASK = comp_board_bitmask{
	start_byte: 3,
	masks:      []uint8{BITMASK_6, BITMASK_7, BITMASK_0, BITMASK_1, BITMASK_2, BITMASK_3},
}

var COLUMN_6_BITMASK = comp_board_bitmask{
	start_byte: 4,
	masks:      []uint8{BITMASK_4, BITMASK_5, BITMASK_6, BITMASK_7, BITMASK_0, BITMASK_1},
}

var HEIGHT_BITMASKS = [7]comp_board_bitmask{HEIGHT_0_BITMASK,
	HEIGHT_1_BITMASK,
	HEIGHT_2_BITMASK,
	HEIGHT_3_BITMASK,
	HEIGHT_4_BITMASK,
	HEIGHT_5_BITMASK,
	HEIGHT_6_BITMASK,
}

var COLUMN_BITMASKS = [7]comp_board_bitmask{COLUMN_0_BITMASK,
	COLUMN_1_BITMASK,
	COLUMN_2_BITMASK,
	COLUMN_3_BITMASK,
	COLUMN_4_BITMASK,
	COLUMN_5_BITMASK,
	COLUMN_6_BITMASK,
}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

var HEIGHT_0_BITMASK = comp_board_bitmask{
	start_byte: 5,
	masks:      []uint8{BITMASK_2, BITMASK_3, BITMASK_4},
}

var HEIGHT_1_BITMASK = comp_board_bitmask{
	start_byte: 5,
	masks:      []uint8{BITMASK_5, BITMASK_6, BITMASK_7},
}

var HEIGHT_2_BITMASK = comp_board_bitmask{
	start_byte: 6,
	masks:      []uint8{BITMASK_0, BITMASK_1, BITMASK_2},
}

var HEIGHT_3_BITMASK = comp_board_bitmask{
	start_byte: 6,
	masks:      []uint8{BITMASK_3, BITMASK_4, BITMASK_5},
}

var HEIGHT_4_BITMASK = comp_board_bitmask{
	start_byte: 7,
	masks:      []uint8{BITMASK_6, BITMASK_7, BITMASK_0},
}

var HEIGHT_5_BITMASK = comp_board_bitmask{
	start_byte: 7,
	masks:      []uint8{BITMASK_1, BITMASK_2, BITMASK_3},
}

var HEIGHT_6_BITMASK = comp_board_bitmask{
	start_byte: 7,
	masks:      []uint8{BITMASK_4, BITMASK_5, BITMASK_6},
}

var TURN_BITMASK = comp_board_bitmask{
	start_byte: 7,
	masks:      []uint8{BITMASK_7},
}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type compressed_board struct {
	raw []uint8
}
type basic_board struct {
	board [][]uint8
	turn  uint8
}
type rich_board struct {
	board  [][]uint8
	turn   uint8
	winner uint8
	tie    bool
}

func (cb compressed_board) check_spot(column_ind uint8, height_int uint8) (output uint8, err error) {
	if cb.read_height(column_ind) < height_int {
		err = errors.New("invalid spot: Spot exists above highest piece in column")
		return 0, err
	}
	if COLUMN_BITMASKS[column_ind].masks[height_int] > COLUMN_BITMASKS[column_ind].masks[0] {
		output = cb.raw[COLUMN_BITMASKS[column_ind].start_byte+1] & COLUMN_BITMASKS[column_ind].masks[height_int]
	}

	output = cb.raw[COLUMN_BITMASKS[column_ind].start_byte] & COLUMN_BITMASKS[column_ind].masks[height_int]

	return output, nil
}

func (cb compressed_board) read_height(column_ind uint8) (height uint8) {
	var byte_offset uint8
	for mask_ind, mask := range HEIGHT_BITMASKS[column_ind].masks {
		if cb.raw[HEIGHT_BITMASKS[column_ind].start_byte+byte_offset]&mask > 0 {
			height = height & (1 << (2 - mask_ind))
		}
		if mask == BITMASK_7 {
			byte_offset = 1
		}
	}
	return height
}

func (cb compressed_board) check_spot_distances(column_ind uint8, height_ind uint8) (
	up_right_valid bool,
	down_right_valid bool,
	down_left_valid bool,
	up_left_valid bool, err error) {

	if cb.read_height(column_ind) < height_ind {
		err = errors.New("spot is out of height")
		return
	}
	var misses_left bool = column_ind >= 3
	var misses_right bool = column_ind <= 3

	var misses_floor bool = height_ind >= 3
	var misses_ceiling bool = height_ind <= 2

	up_right_valid = misses_right && misses_ceiling
	down_right_valid = misses_right && misses_floor
	up_left_valid = misses_left && misses_ceiling
	down_left_valid = misses_left && misses_floor

	for col := uint8(column_ind) + 1; col < uint8(MAX_COLUMNS); col++ {
	}
	return
}

func (cb compressed_board) check_diags_from_end(column_ind uint8, height_ind uint8) (is_win bool, err error) {
	win_piece, err := cb.check_spot(column_ind, height_ind)
	if err != nil {
		return false, err
	}

	up_right_valid, down_right_valid, down_left_valid, up_left_valid, err := cb.check_spot_distances(column_ind, height_ind)
	if err != nil {
	}

	var increment_funcs [4]func(*uint8, *uint8)

	if up_right_valid {
		increment_funcs[0] = func(height *uint8, column *uint8) { *height += 1; *column += 1 }
	}
	if down_right_valid {
		increment_funcs[1] = func(height *uint8, column *uint8) { *height -= 1; *column += 1 }
	}
	if up_left_valid {
		increment_funcs[2] = func(height *uint8, column *uint8) { *height += 1; *column -= 1 }
	}
	if down_left_valid {
		increment_funcs[3] = func(height *uint8, column *uint8) { *height -= 1; *column -= 1 }
	}

	var working_col_ind uint8
	var working_height_ind uint8
	var working_piece uint8
	for _, inc_func := range increment_funcs {
		working_col_ind = column_ind
		working_height_ind = height_ind
		for i := 0; i < 3; i++ {
			inc_func(&working_col_ind, &working_height_ind)
			//TODO: CHECK FOR INVALID HEIGHTS
			working_piece, err = cb.check_spot(working_col_ind, working_height_ind)
			if err != nil {
				return false, errors.New(fmt.Sprintf("check_diags_from_end recieved error from cb.check_spot: {%s}", err))
			}
			if working_piece != win_piece {
				continue
			}

		}
		if working_piece != win_piece {
			continue
		}
		return true, nil
	}
	return false, nil
}

// func (cb compressed_board) to_basic_board() (bb basic_board) {
// 	board := make([][]uint8, MAX_COLUMNS+1)
// 	for col_ind, mask_group := range COLUMN_BITMASKS {
// 		column := make([]uint8, MAX_HEIGHT+1)
// 		for _, mask := mask_group.masks {
//
// 		}
// 	}
// }

func (cb *compressed_board) apply_turn_in_place(turn_column uint8) (err error) {

	if turn_column > MAX_COLUMNS {
		return errors.New(fmt.Sprintf("Invalid input. There is no column %i. Max Column is %i", turn_column, MAX_COLUMNS))
	}

	// add height
	var height uint8
	var byte_offset uint8
	cb.read_height(turn_column)

	if height >= MAX_HEIGHT {
		return errors.New(fmt.Sprintf("Illegal Move. Column %i is already full.", turn_column))
	}
	height_diff := height ^ (height + 1)
	var bit_flips [3]bool = [3]bool{height_diff > 3, height_diff%4 > 1, height_diff%2 != 0}

	for mask_ind, mask := range HEIGHT_BITMASKS[turn_column].masks {
		if bit_flips[mask_ind] {
			cb.raw[HEIGHT_BITMASKS[turn_column].start_byte+byte_offset] = cb.raw[HEIGHT_BITMASKS[turn_column].start_byte+byte_offset] ^ mask
		}

		if mask == BITMASK_7 {
			byte_offset = 1
		}
	}

	// check if we need a piece
	var current_turn uint8 = TURN_BITMASK.masks[0] & cb.raw[TURN_BITMASK.start_byte]
	if current_turn != 0 {
		//add piece if necessary
		if COLUMN_BITMASKS[turn_column].masks[height] > COLUMN_BITMASKS[turn_column].masks[0] {
			cb.raw[COLUMN_BITMASKS[turn_column].start_byte] = cb.raw[COLUMN_BITMASKS[turn_column].start_byte] & COLUMN_BITMASKS[turn_column].masks[height]
		} else {
			cb.raw[COLUMN_BITMASKS[turn_column].start_byte+1] = cb.raw[COLUMN_BITMASKS[turn_column].start_byte+1] & COLUMN_BITMASKS[turn_column].masks[height]
		}
	}

	// change turn
	cb.raw[TURN_BITMASK.start_byte] = TURN_BITMASK.masks[0] ^ cb.raw[TURN_BITMASK.start_byte]
	return nil
}

func (cb compressed_board) apply_turn_copy(turn_column uint8) compressed_board {
	cb.apply_turn_in_place(turn_column)
	return cb
}
