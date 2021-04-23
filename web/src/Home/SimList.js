import React from 'react';
import {Link} from 'react-router-dom';

const SimList = ({sims}) => {
    return sims.map(sim => {
        return(
            <div key={sim.ID}>
                <Link to={'/simulation/' + sim.ID}>{sim.Name}</Link>
                <div>{sim.Description}</div>
            </div>
        );
    })
}

export default SimList