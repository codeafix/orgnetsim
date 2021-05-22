import React, { useEffect, useRef, useState } from 'react';
import { select, forceSimulation, forceLink, forceManyBody, forceCenter, forceX, forceY, drag, zoom } from 'd3';
import Color from './Color';
import API from '../api';
import Spinner from 'react-bootstrap/Spinner';
import ToggleButton from 'react-bootstrap/ToggleButton';
import ButtonGroup from 'react-bootstrap/ButtonGroup';

const NetworkGraph = (props) => {
    const graph = useRef(null);
    const [loading, setloading] = useState(false);
    const [layout, setlayout] = useState(false);
    const [step, setstep] = useState();
    const [run, setrun] = useState({stopped:true, sim:null});

    const runsim = (enable) => {
        setlayout(enable);
        run.stopped = !enable;
        if(enable){
            run.sim.restart();
        }else{
            run.sim.stop();
            savegraph();
        }
    };

    const savegraph = () => {
        if(!run.sim || !step) return;

        var graphnodes = run.sim.nodes();
        var netnodes = step.network.nodes;
        for(var i = 0; i < graphnodes.length; i++) {
            var uinode = graphnodes[i];
            var netnode = netnodes[i];
            netnode.fx = uinode.x;
            netnode.fy = uinode.y
        }
        var links = [];
        var graphlinks = run.sim.force("link").links();
        for(var i = 0; i < graphlinks.length; i++) {
            var link = {
                source: graphlinks[i].source.id,
                target: graphlinks[i].target.id,
            };
            links.push(link);
        }
        step.network.links = links;
        var response = API.updateStep(step);
        response.then(updatedstep => setstep(updatedstep)
        ).catch(err => {
            console.log(err);
        });
    };

    const createGraph = (network) => {
        const margin = {top: 0, right: 0, bottom: 0, left: 0},
            vwidth = graph.current.parentElement.offsetWidth,
            cw = vwidth - margin.left - margin.right,
            vheight = Math.round(cw/1.6),
            ch = vheight - margin.top - margin.bottom;
        
        const width = 1024, height = 768;

        const resize = () => {
            var w = graph.current.parentElement.offsetWidth,
                h = Math.round((w -  margin.left - margin.right)/1.6);
            select(graph.current).attr('width', w)
                .attr('height', h);
        };

        select(window).on(
            'resize.' + select(graph.current.parentElement).attr('id'), 
            resize
        );
        
        const svg = select(graph.current);

        svg.attr('viewBox', `0 0 ${cw} ${ch}`)
            .attr('preserveAspectRatio', 'xMinYMid')
            .call(resize);
        
        const simulation = forceSimulation()
            .force("link", forceLink().id(function(d) { return d.id; }))
            .force("charge", forceManyBody())
            .force("center", forceCenter(width / 2, height / 2))
            .force("xAxis", forceX().strength(0.01).x((width)/2))
            .force("yAxis", forceY().strength(0.01).y((height)/2));

        const networkGraph = svg.append('svg:g').attr('class','grpParent');

        const zoomed = (event) => {
            networkGraph.attr("transform", event.transform);
        }

        const z = zoom()
            .scaleExtent([-40, 40])
            .translateExtent([[-2*width, -2*height], [4*width, 4*height]])
            .on("zoom", zoomed);

        svg.call(z);

        const its = 2000, chgits = 2000;

        const links = networkGraph.append("g")
            .attr("class", "links")
            .selectAll("line")
            .data(network.links)
            .enter().append("line")
            .style("stroke", "LightGray")
            .style('stroke-width', function(d) { return 5*d.strength/its;} );
        
        const nodes = networkGraph.append("g")
            .attr("class", "nodes")
            .selectAll("circle")
            .data(network.nodes)
            .enter().append("circle")
            .attr("r", function(d,i){return 5 + 2.5*d.change/chgits;})
            .style("fill", function(d){ return Color.cssColorFromVal(d.color); })
            .call(drag()
                .on("start", dragstarted)
                .on("drag", dragged)
                .on("end", dragended));
      
        nodes.append("title")
            .text(function(d) { return d.id; });
      
        const ticked = () => {
            links
                .attr("x1", function(d) { return d.source.x; })
                .attr("y1", function(d) { return d.source.y; })
                .attr("x2", function(d) { return d.target.x; })
                .attr("y2", function(d) { return d.target.y; });
        
            nodes
                .attr("cx", function(d) { 
                    return d.x = Math.max(-2*width+5, Math.min(4*width - 5, d.x)); })
                .attr("cy", function(d) { 
                    return d.y = Math.max(-2*height+5, Math.min(4*height - 5, d.y)); });
            };

        simulation
            .nodes(network.nodes)
            .on("tick", ticked);
      
        simulation.force("link")
            .links(network.links);
        
        run.sim = simulation;
        runsim(false);
        simulation.tick();
        ticked();
        
        function dragstarted(event, d) {
            if(run.stopped) return;
            if (!event.active) simulation.alphaTarget(0.3).restart();
            d.fx = d.x;
            d.fy = d.y;
        };
            
        function dragged(event, d) {
            if(run.stopped) return;
            d.fx = event.x;
            d.fy = event.y;
        };
            
        function dragended(event, d) {
            if(run.stopped) return;
            if (!event.active) simulation.alphaTarget(0);
            d.fx = null;
            d.fy = null;
        };
    }

    useEffect(() => {
        if (!props.sim['id']) {
            return;
        }
        const steps = props.sim.steps || [];
        const sl = steps.length;
        if (sl == 0){
            return;
        }
        setloading(true);
        select(graph.current).selectAll("*").remove();
        
        const lastStep = props.sim.steps[sl-1];
        API.getStep(lastStep).then(step => {
            setstep(step);
            createGraph(step.network);
            setloading(false);
        }).catch(err => {
            console.log(err);
            setloading(false);
        });
    },[props.sim]);

    return <div>
            {loading && <Spinner animation="border" variant="info" />}
            <svg class="mb-3" ref={graph}/>
            <ButtonGroup size="sm" toggle className="btn btn-primary float-right" >
                <ToggleButton type="checkbox" checked={layout} onChange={(e) => runsim(e.currentTarget.checked)}>{layout ? 'Click to save' : 'Click to layout'}</ToggleButton>
            </ButtonGroup>
        </div>
}

export default NetworkGraph;