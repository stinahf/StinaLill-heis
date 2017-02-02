package elev

import (
    def "config"
    "errors"
    "log"
)

const MOTOR_SPEED 2800 //TODO - Move to config file


var lamp_channel_matrix[def.N_FLOORS][def.N_BUTTONS]int{
    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}


var button_channel_matrix[def.N_FLOORS][def.N_BUTTONS]int{
    {BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}


func Elev_init() int {
    if !io_init() {
        return -1
    }

    for f := 0; f < def.N_FLOORS; f++ {
        if f != 0 {
            Elev_set_button_lamp(f, def.BUTTON_UP, false)
        }
        if f != def.N_FLOORS-1 {
            Elev_set_button_lamp(f, def.BUTTON_DOWN, false)
        }
        Elev_set_button_lamp(f, def.BUTTON_INSIDE, false)
    }

    Elev_set_stop_lamp(false)
    Elev_set_door_lamp(false)
}


func Elev_set_motor_direction(dirn int) {
    if (dirn == 0){
        io_write_analog(MOTOR, 0)
    } else if dirn > 0 {
        io_clear_bit(MOTORDIR);
        io_write_analog(MOTOR, MOTOR_SPEED);
    } else if (dirn < 0) {
        io_set_bit(MOTORDIR);
        io_write_analog(MOTOR, MOTOR_SPEED);
    }
}


func Elev_set_button_lamp(button int, floor int, value bool) {
    if floor < 0 || floor >= def.N_FLOORS {
        log.Printf("Error: The floor is out of range", floor)
        return
    }
    if button == def.BUTTON_UP && floor == def.N_FLOORS-1 {
        log.Println("You are already at the top")
        return
    }
    if button == def.BUTTON_DOWN && floor == 0 {
        log.Println("You are already at the bottom")
        return
    }
    if button != def.BUTTON_UP && button != def.BUTTON_DOWN && button != def.BUTTON_INSIDE {
        log.Printf("Invalid button %d\n", button)
        return
    }

    if value {
        io_set_bit(lamp_channel_matrix[floor][button])
    } else {
        io_clear_bit(lamp_channel_matrix[floor][button])
    }
}


func Elev_set_floor_indicator(floor int) {
    if floor < 0 || floor >= def.N_FLOORS {
        log.Printf("The floor %d is out of range! \n", floor)
        return
    }

    // Binary encoding. One light must always be on.
    if floor & 0x02 > 0 {
        io_set_bit(LIGHT_FLOOR_IND1);
    } else {
        io_clear_bit(LIGHT_FLOOR_IND1);
    }    

    if floor & 0x01 > 0 {
        io_set_bit(LIGHT_FLOOR_IND2);
    } else {
        io_clear_bit(LIGHT_FLOOR_IND2);
    }    
}


func Elev_set_door_open_lamp(value bool) {
    if value {
        io_set_bit(LIGHT_DOOR_OPEN);
    } else {
        io_clear_bit(LIGHT_DOOR_OPEN);
    }
}


func Elev_set_stop_lamp(value bool) {
    if value {
        io_set_bit(LIGHT_STOP);
    } else {
        io_clear_bit(LIGHT_STOP);
    }
}



func Elev_get_button_signal(button int, floor int) bool {
    if floor < 0 || floor >= def.N_FLOORS {
        log.Printf("The floor %d is out of range \n", floor)
        return false
    }
    if button < 0 || button >= def.N_BUTTONS {
        log.Printf("Button %d is out of range \n", floor)
        return false
    }
    if button == def.BUTTON_UP && floor == def.N_FLOORS-1 {
        log.Println("You are already on the top")
        return
    }
    if button == def.BUTTON_DOWN && floor == 0 {
        log.Println("You are already at the buttom")
        return
    }
    if io_read_bit(button_channel_matrix[floor][button]) {
        return true
    } else {
        return false
    }
}


func Elev_get_floor_sensor_signal() int {
    if io_read_bit(SENSOR_FLOOR1) {
        return 0
    } else if io_read_bit(SENSOR_FLOOR2) {
        return 1
    } else if io_read_bit(SENSOR_FLOOR3) {
        return 2
    } else if io_read_bit(SENSOR_FLOOR4) {
        return 3
    } else {
        return -1
    }
}


func Elev_get_stop_signal() bool{
    return io_read_bit(STOP);
}


func Elev_get_obstruction_signal() bool{
    return io_read_bit(OBSTRUCTION);
}

