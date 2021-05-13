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
 
    colorValArray: [
        0,
        1,
        2,
        3,
        4,
        5,
        6
    ],
    
    //return the color name used in the UI for the given color enum
    colorFromVal: function(color) {
        if (color >= 0 && color < this.nameArray.length){
            return this.nameArray[color];
        }
        return "Invalid Color";
    },
    //return the color style name used to fill objects in a d3 chart
    cssColorFromVal: function(color) {
        if (color >= 0 && color < this.nameArray.length){
            return this.cssColorArray[color];
        }
        return this.cssColorArray[0];
    },

    //return a slice of the UI color names
    colorNameSlice: function(maxColors) {
        return this.nameArray.slice(0,maxColors);
    },

    //return a slice of the UI color values
    cssValSlice: function(maxColors) {
        return this.cssColorArray.slice(0,maxColors);
    },

    //return a slice of the UI color values
    colorValSlice: function(maxColors) {
        return this.colorValArray.slice(0,maxColors);
    }
}

export default Color