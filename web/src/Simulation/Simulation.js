import React from 'react';
import API from '../api';
import {Link} from 'react-router-dom';

const Simulation = (props) => {
    const sim = API.get(
        parseInt(props.match.params.number, 10)
    )
    if (!sim) {
        return <div>Sorry, but the simulation was not found</div>
    }
    return (
        <div>
            <h1>{sim.Name}</h1>
            <p>{sim.Description}</p>
            <Link to='/'>Home</Link>
        </div>
    )
}

export default Simulation