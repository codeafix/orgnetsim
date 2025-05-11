import React, { useEffect, useRef, useState } from 'react';
import { select, forceSimulation, forceLink, forceManyBody, forceCenter, forceX, forceY, drag, zoom, Simulation, SimulationNodeDatum, SimulationLinkDatum, ForceLink, Selection } from 'd3';
import Color from './Color';
import API from '../API/api';
import Spinner from 'react-bootstrap/Spinner';
import {Button} from 'react-bootstrap';
import {indexBy} from 'underscore';
import { SimInfo } from '../API/SimInfo';
import { Step } from '../API/Step';
import type { Network, Link, AgentState } from '../API/Network';

type NetworkGraphProps = {
    sim:SimInfo;
    steps:Array<Step>;
}

type Run = {
    stopped: boolean;
    sim?: Simulation<SimulationNodeDatum, SimulationLinkDatum<SimulationNodeDatum>>;
    nodes?: Selection<SVGCircleElement, SimulationNodeDatum, SVGGElement, unknown>;
}

const NetworkGraph = (props:NetworkGraphProps) => {
    const graph = useRef<SVGSVGElement>(null);
    const [loading, setloading] = useState<boolean>(false);
    const [layout, setlayout] = useState<boolean>(false);
    const [step, setstep] = useState<Step>();
    const [networkData, setNetworkData] = useState<Network | null>(null);
    const [run,] = useState<Run>({stopped:true});
    const [steps, setsteps] = useState<Array<Step>>([]);
    const [play,] = useState({playing: false});
    const [isrunning, setisrunning] = useState<boolean>(false);
    const [runcount, setruncount] = useState<number>(0);

    const getrunning = () => {
        return play.playing;
    };

    const setrunning = (running:boolean) => {
        setisrunning(running);
        return play.playing = running;
    };

    const runsim = (enable:boolean) => {
        setlayout(enable);
        run.stopped = !enable;
        if(!run.sim) return;
        if(enable){
            run.sim.restart();
        }else{
            run.sim.stop();
            savegraph();
        }
    };

    const savegraph = () => {
        if(!run.sim || !step || !networkData) return;

        var graphnodes = run.sim.nodes();
        var netnodes = networkData.nodes;
        for(var i = 0; i < graphnodes.length; i++) {
            var uinode = graphnodes[i];
            var netnode = netnodes[i];
            netnode.fx = uinode.x || 0;
            netnode.fy = uinode.y || 0
        }
        var links = [];
        var force = run.sim.force<ForceLink<SimulationNodeDatum, SimulationLinkDatum<SimulationNodeDatum>>>("link");
        if(!force) return;
        var graphlinks = force.links();
        for(var j = 0; j < graphlinks.length; j++) {
            var link:Link = {
                source: (graphlinks[j].source as any).id,
                target: (graphlinks[j].target as any).id,
                strength: 0,
                length: 0
            };
            links.push(link);
        }
        
        const updatedNetworkData = {...networkData, links};
        try {
            const response = API.updateStepNetwork(step.parent, step.id, updatedNetworkData);
            response.then(updatedNetwork => {
                setNetworkData(updatedNetwork);
            });
        } catch (err) {
            console.log(err);
        }
    };

    const hold = (ms:number) => {
        return new Promise(resolve => setTimeout(resolve, ms));
    };

    const playsteps = async () => {
        if(getrunning()) {
            setrunning(false);
            return;
        }
        setrunning(true);
        if(run.sim){
            var graphnodes = run.sim.nodes();
            var graphnodesbyid = indexBy(graphnodes, (n:AgentState) => n.id);
            var stepcount = steps.length;
            var waittime = 10000 / stepcount;
            for(var i = 0; i < stepcount; i++) {
                if(!getrunning()) return;
                setruncount(stepcount - i);
                var s = steps[i];
                await renderstep(s, graphnodesbyid, waittime);
            }
        }
        setrunning(false);
    };

    const renderstep = async (step:Step, graphnodesbyid:any, waittime:number) => {
        if(!run.sim || !run.nodes) return;
        var netnodes = await API.getStepAgents(step.parent, step.id) || [];
        for(var i = 0; i < netnodes.length; i++) {
            var node = netnodes[i];
            var graphnode = graphnodesbyid[node.id];
            graphnode.color = node.color;
        }
        run.sim.tick();
        run.nodes.style("fill", function(d:any){ return Color.cssColorFromVal(d.color); });
        await hold(waittime);
    };

    const createGraph = (network: Network | null) => {
        if (!network) return;
        if(!graph.current) return;
        const parentElement = graph.current.parentElement;
        if(!parentElement) return;
        const margin = {top: 0, right: 0, bottom: 0, left: 0},
            vwidth = parentElement.offsetWidth,
            cw = vwidth - margin.left - margin.right,
            vheight = Math.round(cw/1.6),
            ch = vheight - margin.top - margin.bottom;
        
        const width = 1024, height = 768;
        
        const resize = () => {
            const c = graph.current;
            if(!c) return;
            if(!c.parentElement) return;
            var w = c.parentElement.offsetWidth,
                h = Math.round((w -  margin.left - margin.right)/1.6);
            select(c).attr('width', w)
                .attr('height', h);
        };

        select(window).on(
            'resize.' + select(parentElement).attr('id'), 
            resize
        );
        
        const svg = select(graph.current);

        svg.attr('viewBox', `0 0 ${cw} ${ch}`)
            .attr('preserveAspectRatio', 'xMinYMid')
            .call(resize);
        
        const simulation = forceSimulation<SimulationNodeDatum, SimulationLinkDatum<SimulationNodeDatum>>()
            .force("link", forceLink().id(function(d:any) { return d.id; }))
            .force("charge", forceManyBody())
            .force("center", forceCenter(width / 2, height / 2))
            .force("xAxis", forceX().strength(0.01).x((width)/2))
            .force("yAxis", forceY().strength(0.01).y((height)/2));

        const networkGraph = svg.append('svg:g').attr('class','grpParent');

        const zoomed = (event:any) => {
            networkGraph.attr("transform", event.transform);
        }

        const z = zoom<SVGSVGElement, unknown>()
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
            .style('stroke-width', function(d) { return d.strength ? 10*(d.strength)/its : 1; });
        
        const nodes = networkGraph.append("g")
            .attr("class", "nodes")
            .selectAll("circle")
            .data(network.nodes)
            .enter().append("circle")
            .attr("r", function(d,i){return 5 + 5*d.change/chgits;})
            .style("fill", function(d){ return Color.cssColorFromVal(d.color); })
            .call(drag<SVGCircleElement, AgentState>()
                .on("start", dragstarted)
                .on("drag", dragged)
                .on("end", dragended));
      
        nodes.append("title")
            .text(function(d) { return d.name; });
      
        const ticked = () => {
            links
                .attr("x1", function(d:any) { return d.source.x; })
                .attr("y1", function(d:any) { return d.source.y; })
                .attr("x2", function(d:any) { return d.target.x; })
                .attr("y2", function(d:any) { return d.target.y; });
        
            nodes
                .attr("cx", function(d) { 
                    return d.x = Math.max(-2*width+5, Math.min(4*width - 5, d.x||0)); })
                .attr("cy", function(d) { 
                    return d.y = Math.max(-2*height+5, Math.min(4*height - 5, d.y||0)); });
            };

        simulation
            .nodes(network.nodes)
            .on("tick", ticked);
      
        var force = simulation.force<ForceLink<SimulationNodeDatum, SimulationLinkDatum<SimulationNodeDatum>>>("link");
        if (!force) return;
        force.links(network.links);
        
        run.sim = simulation;
        run.nodes = nodes as any;
        runsim(false);
        simulation.tick();
        ticked();
        
        function dragstarted(event:any, d:AgentState) {
    if(run.stopped) return;
    if (!event.active) simulation.alphaTarget(0.3).restart();
    d.fx = d.x;
    d.fy = d.y;
};
            
function dragged(event:any, d:AgentState) {
    if(run.stopped) return;
    d.fx = event.x;
    d.fy = event.y;
};
            
function dragended(event:any, d:AgentState) {
    if(run.stopped) return;
    if (!event.active) simulation.alphaTarget(0);
    d.fx = undefined;
    d.fy = undefined;
};
}

useEffect(() => {
        setloading(true);
        if (!props.sim['id']) {
            setloading(false);
            return;
        }
        const steplist = props.steps || [];
        setsteps(steplist);

        if (props.sim.steps.length === 0) setloading(false);
        if (steplist.length === 0) return;

        select(graph.current).selectAll("*").remove();

        const lastStep = steplist[steplist.length-1];
        setstep(lastStep);

        try {
            // Fetch network data for the last step
            const response = API.getStepNetwork(props.sim.id, lastStep.id);
            response.then(networkData => {
                setNetworkData(networkData);
                createGraph(networkData);
            });
        } catch (error) {
            console.error('Failed to fetch network data', error);
        } finally {
            setloading(false);
        }
    },[props.sim, props.steps]);//eslint-disable-line react-hooks/exhaustive-deps

    return <div id="graph-container">
            {loading && <Spinner animation="border" variant="info" />}
            <svg className="mb-3" ref={graph}/>
            <Button size="sm" className="btn btn-primary float-right" onClick={(e) => playsteps()} disabled={steps.length < 2 || layout}>{isrunning ? runcount : 'Play'}</Button>
            <Button size="sm" className="btn btn-primary float-right" onClick={(e) => runsim(!layout)} disabled={steps.length > 1}>{layout ? 'Save layout' : 'Layout'}</Button>
        </div>
}

export default NetworkGraph;