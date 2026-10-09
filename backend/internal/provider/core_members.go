package provider

import (
	"github.com/Ankumeah/JSBEE/backend/internal/frontend/assets"
)

type member struct {
	Name        string
	ImagePath   string
	Description string
}

var CoreMembers = map[string]member{
	"Founder & Executive Lead": {
		Name:      "Ayush Anand",
		ImagePath: "/" + assets.AyushAnandImage,
		Description: `Ayush Anand is a student, writer and researcher with a knack for building things. He cares deeply about business, economics and sustainability, especially about how young people can actually contribute to change instead of just reading about it. He founded the Journal of Sustainable Business in Emerging Economies (JSBEE), a student-led journal that gives young people room to explore sustainable business, economics and emerging economies through their own research.

Outside JSBEE he has picked up experience across research, writing, publishing and youth work, from community projects to student-led organisations. What he reads and writes about keeps widening. Business and economics sit next to sociology, psychology and philosophy, with sustainability and social issues running through all of it. He has also researched sports, social integration and youth development, and he keeps looking for ways to get more students into research and academic programmes.

He writes too. His book Destined Hearts explores young emotions and relationships, and beyond fiction he turns out poetry, essays and reflective pieces. A lot of his writing circles the same question: what are the ideas and experiences people do not usually say out loud?

What ties his projects together is a simple idea. School only takes you so far, so he builds spaces where students can try research, leadership, communication and entrepreneurship for themselves and work on real problems. The bigger goal is student-led platforms where young people find opportunities, meet each other and make work they actually own.

Day to day at JSBEE that means laying the foundation: growing the team, shaping its programmes, lining up collaborations and setting the journal's long-term academic direction.`,
	},
	"Co-founder & Managing Director": {
		Name:        "Anirban",
		ImagePath:   "/" + assets.AnirbanImage,
		Description: `Anirban is drawn to academics and research. You will find him at Model United Nations conferences and student competitions, and he is known as a confident speaker. He has also published a book.`,
	},
	"Co-founder & Director of Research and Academics": {
		Name:        "Aryan Kumar",
		ImagePath:   "/" + assets.AryanKumarImage,
		Description: `Aryan co-founded JSBEE and runs its academic side. He sets the bar for research, guides contributing authors and keeps the journal's academic backbone in shape.`,
	},
	"Content & Communication Lead": {
		Name:      "Shivam Kumar",
		ImagePath: "/" + assets.ShivamKumarImage,
		Description: `Shivam has been building a real audience on LinkedIn: 500+ followers and 50+ posts on biotechnology, careers, student life and the lessons picked up along the way. In one consistent run his posts pulled around 18K impressions.

That volume taught him what actually works. Hooks that grab, complicated ideas made simple and stories people see themselves in. He keeps trying new formats to learn what readers respond to and how to say things plainly without dumbing them down.`,
	},
	"PR and Partnership": {
		Name:        "Anushka Singh",
		ImagePath:   "/" + assets.AnushkaSinghImage,
		Description: `Anushka is a Class 12 Commerce student at The Pentecostal Assembly School in Bokaro who likes business, economics, finance and research. She topped her class in Class 11 with 92% and Rank 1. Beyond school she was Best Delegate at UNGA, won Youth Parliament, picked up a Chess Champion Award and has competed in inter-school Science and Commerce Olympiads. She writes as well, with newspaper articles and self-composed poems in magazines to her name, which is where her research and communication chops come from. Put together it adds up to strong leadership, public speaking, analytical thinking, presentation and teamwork skills.`,
	},
	"International and City Network Lead": {
		Name:        "Divya",
		ImagePath:   "/" + assets.DivyaImage,
		Description: `Divya is a Class 11 PCMB student with a packed academic and extracurricular record. She reads widely and loves writing, literature, history, politics, international relations, diplomacy and public speaking. Curious and analytical, she likes digging into ideas, questioning assumptions and looking at problems from every side. That shows in her skills: critical thinking, analytical reasoning, communication, leadership, creativity, research, problem-solving and independent learning, especially anything with discussion, argument and negotiation. Her shelf is full. Best Delegate at a Model UN in the UNSC, the EI ASSET Prize three times, Zonal Rank 1 in the International English Olympiad and 92.8% in AISSE Class 10, plus school-level awards in elocution, science demonstrations, academic merit, chess and plenty of inter-school contests.`,
	},
	"Finance and Administration": {
		Name:      "Omkar Prasad",
		ImagePath: "/" + assets.OmkarPrasadImage,
		Description: `Omkar Prasad is doing Class 12 (CBSE) in Commerce: Accountancy, Economics, Business Studies, Applied Mathematics and English. Finance, economics, financial markets and research are where his interests lie. He already understands the basics of financial markets and quantitative finance and is working on his own research into financial market volatility.

On the organising side he has run an online community of 70+ members and put together events with 100+ participants. Planning, coordination, documentation, budgeting and record-keeping are his thing, backed by solid leadership, communication and team-management skills.`,
	},
	"Technology and Digital Lead": {
		Name:      "Ritum Prabhat",
		ImagePath: "/" + assets.RitumPrabhatImage,
		Description: `Young entrepreneur, technology innovator and student leader.

Ritum Prabhat is a young entrepreneur and tech enthusiast who works across entrepreneurship, AI, Python, web development, innovation and leadership. He is Managing Director of the Reshvah Initiative, where he looks after digital strategy, the website, research and organisational growth.

His projects include VisionStep, an assistive-tech concept that got recognised at a school STEM fair. He has also been part of MUNs, entrepreneurship events and student leadership initiatives.

Core skills: entrepreneurship, business strategy, AI and Python, web development, innovation, leadership, public speaking, research and strategic thinking.`,
	},
	"Marketing and Outreach Lead": {
		Name:        "Shanmugasundaram",
		ImagePath:   "/" + assets.ShanmugasundaramImage,
		Description: `Shanmugasundaram runs marketing and outreach at JSBEE. He has a gift for winning people over and currently handles marketing and sales for Alaada, so growing JSBEE's reach comes naturally to him.`,
	},
	"Managing Team": {
		Name:        "Kaashvi Kumar",
		ImagePath:   "/" + assets.KaashviKumarImage,
		Description: `Kaashvi Kumar founded Project Cangro, a youth-led advocacy initiative working on reproductive rights through awareness, education and open conversation around taboo topics. Outside JSBEE she has also led social media at Think Economics, coordinated design work at Aletheia and volunteered her design skills with Lumina Aid and ByteMinds. Along the way she has picked up certifications across data science, Microsoft Excel, business analytics, operations management, consumer awareness, Google tools and international relations, won a Judge's Choice Award in a poetry competition and collected Best Delegate and High Commendation awards at Model UNs.`,
	},
	"Programmes & Operations Lead": {
		Name:        "Tanishka Gupta",
		ImagePath:   "/" + assets.TanishkaGuptaImage,
		Description: `Tanishka Gupta is a Grade 11 Commerce student who enjoys event planning, teamwork and community building. She has organised school events, taken part in student council and AFS workshops, and brings adaptability, strong collaboration skills and a willingness to learn. At JSBEE she wants to help plan events, coordinate logistics and build an active, inclusive student community while growing her own leadership and project management skills.`,
	},
}

var CoreMemberOrder = []string{
	"Founder & Executive Lead",
	"Co-founder & Managing Director",
	"Co-founder & Director of Research and Academics",
	"Content & Communication Lead",
	"PR and Partnership",
	"International and City Network Lead",
	"Finance and Administration",
	"Technology and Digital Lead",
	"Marketing and Outreach Lead",
	"Managing Team",
	"Programmes & Operations Lead",
}
