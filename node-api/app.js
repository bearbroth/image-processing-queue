import express from 'express';
import { publishJob, initializeQueue } from './queue.js';
import redis from './redisClient.js';
import { v4 } from 'uuid';
const app = express();
const PORT = 3000

app.use(express.json());
app.use(express.urlencoded({extended : true}));

app.listen(PORT, ()=>{
    console.log('Listening...');
})

initializeQueue();


app.post('/process',async(req,res) =>{
    const image = req.body.image;
    const jobID = v4();
    await redis.set(`job:${jobID}`,'queued')
    publishJob({jobID, image});

    res.json({
        jobID : jobID,
        status : 'queued'
    })
})

app.get("/status/:id", async(req,res)=>{
    const jobID = req.params.id;
    const redisKey = `job:${jobID}`;
    const value = await redis.get(redisKey); // value will be queued, processing, failed etc

    if ( !value ){
        res.sendStatus(404);
    }
    else{
        res.json({
            job : jobID,
            status : value
        })
    }
})