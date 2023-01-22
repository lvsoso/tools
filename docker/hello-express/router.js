import Router from 'express';


const router = Router();

router.get('/', function(req, res) {
	return res.json({ hello: "world" });
});

router.post('/', function(req, res) {
	return res.json({
		hello: "world",
		input: req.body,
	});
});

export default router;