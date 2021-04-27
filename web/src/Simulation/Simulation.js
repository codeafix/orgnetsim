import React, { useState, useEffect } from 'react';
import API from '../api';
import {Link} from 'react-router-dom';

const Simulation = (props) => {
    const [sim, setSim] = useState({});
    
    useEffect(() => {
        API.get(props.match.params.id).then(response => {
            setSim(response);
        })
      });

    return (
        <div>
            <h1>{sim.name}</h1>
            <p>{sim.description}</p>
            <Link to='/'>Home</Link>
        </div>
    )
}

export default Simulation