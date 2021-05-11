const Color = {
    nameArray: [
        "Grey",
        "Blue",
        "Red",
        "Green",
        "Yellow",
        "Orange",
        "Purple",
    ],

    cssColorArray: [
        "SlateGray",
        "Blue",
        "Red",
        "LawnGreen",
        "Yellow",
        "Orange",
        "Purple"
    ],
    
    colorFromVal: function(color) {
        if (color >= 0 && color < this.nameArray.length){
            return this.nameArray[color];
        }
        return "Invalid Color";
    },

    cssColorFromVal: function(color) {
        if (color >= 0 && color < this.nameArray.length){
            return this.cssColorArray[color];
        }
        return this.cssColorArray[0];
    }
}

export default Color