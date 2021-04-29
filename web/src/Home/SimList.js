import React from 'react';
import {Link} from 'react-router-dom';
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/esm/Button';

const SimList = ({sims}) => {
    return sims.map(sim => {
            return(
                <Card border="info" style={{ width: '18rem' }} key={sim.id}>
                    <Card.Header><Card.Title>{sim.name}</Card.Title></Card.Header>
                    <Card.Body>
                        <Card.Text>{sim.description}</Card.Text>
                    </Card.Body>
                    <Card.Footer>
                        <Link className="btn btn-primary" role="button" to={'/simulation/' + sim.id}>Open</Link>
                    </Card.Footer>
                </Card>
            );
        });
}

export default SimList