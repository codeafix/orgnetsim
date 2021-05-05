const Color = {
    colorFromVal: function(color) {
        switch(color) {
            case 0:
                return "Grey";
            case 1:
                return "Blue";
            case 2:
                return "Red";
            case 3:
                return "Green";
            case 4:
                return "Yellow";
            case 5:
                return "Orange";
            case 6:
                return "Purple";
            default:
                return "Invalid Color";
        }
    }
}

export default Color