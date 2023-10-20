import{eW as e,eX as t,e1 as a,P as n,h as s,eY as i,eZ as r,E as d,ct as o,co as l,cr as m,a2 as c,F as u,aG as f,br as g,w as h,B as p,aH as S,K as O}from"./detailed_results-2fca7066.chunk.js";import{F as v,e as P,I as b,T as y,a2 as I}from"./preset_utils-c298d622.chunk.js";const T={inputs:[]},E=v({fieldName:"aura",label:"Aura",values:[{name:"None",value:e.NoPaladinAura},{name:"Devotion Aura",value:e.DevotionAura},{name:"Retribution Aura",value:e.RetributionAura}]}),w=v({fieldName:"judgement",label:"Judgement",labelTooltip:"Judgement debuff you will use on the target during the encounter.",values:[{name:"None",value:t.NoJudgement},{name:"Wisdom",value:t.JudgementOfWisdom},{name:"Light",value:t.JudgementOfLight}]}),R={items:[{id:40298,enchant:3819,gems:[41401,40012]},{id:44662,gems:[40012]},{id:40573,enchant:3809,gems:[40012]},{id:44005,enchant:3831,gems:[40012]},{id:40569,enchant:3832,gems:[40012,40012]},{id:40332,enchant:1119,gems:[40012,0]},{id:40570,enchant:3604,gems:[40012,0]},{id:40259,gems:[40012]},{id:40572,enchant:3721,gems:[40027,40012]},{id:40592,enchant:3606},{id:40399},{id:40375},{id:44255},{id:37111},{id:40395,enchant:2666},{id:40401,enchant:1128},{id:40705}]},W={items:[{id:46180,enchant:3820,gems:[41401,40094]},{id:45443,gems:[40012]},{id:46182,enchant:3810,gems:[40012]},{id:45486,enchant:3831,gems:[40012]},{id:45445,enchant:3832,gems:[42148,42148,42148]},{id:45460,enchant:1119,gems:[40012,0]},{id:46179,enchant:3604,gems:[40047,0]},{id:45616,gems:[40012,40012,40012]},{id:46181,enchant:3721,gems:[40012,40012]},{id:45537,enchant:3606,gems:[40012,40012]},{id:45614,gems:[45882]},{id:45946,gems:[40012]},{id:46051},{id:37111},{id:46017,enchant:2666},{id:45470,enchant:1128,gems:[40012]},{id:40705}]},A=P("PreRaid",{items:[{id:44949,enchant:3819,gems:[41401,40012]},{id:42647,gems:[42702]},{id:37673,enchant:3809,gems:[40012]},{id:41609,enchant:3831},{id:39629,enchant:3832,gems:[40012,40012]},{id:37788,enchant:1119,gems:[0]},{id:39632,enchant:3604,gems:[40012,0]},{id:40691,gems:[40012,40012]},{id:37362,enchant:3721,gems:[40012,40012]},{id:44202,enchant:3606,gems:[40094]},{id:44283},{id:37694},{id:44255},{id:37111},{id:37169,enchant:2666},{id:40700,enchant:1128},{id:40705}]}),F=P("P1 Preset",R),j=P("P2 Preset",W),B={name:"Standard",data:a.create({talentsString:"50350151020013053100515221-50023131203",glyphs:{major1:n.GlyphOfHolyLight,major2:n.GlyphOfSealOfWisdom,major3:n.GlyphOfBeaconOfLight,minor2:s.GlyphOfLayOnHands,minor1:s.GlyphOfSenseUndead,minor3:s.GlyphOfBlessingOfKings}})},k=i.create({}),H=r.create({aura:e.DevotionAura,judgement:t.NoJudgement}),M=d.create({defaultPotion:o.RunicManaPotion,flask:l.FlaskOfTheFrostWyrm,food:m.FoodFishFeast});class G extends b{constructor(e,t){super(e,t,{cssClass:"holy-paladin-sim-ui",cssScheme:"paladin",knownIssues:[],epStats:[c.StatIntellect,c.StatSpirit,c.StatSpellPower,c.StatSpellCrit,c.StatSpellHaste,c.StatMP5],epReferenceStat:c.StatSpellPower,displayStats:[c.StatHealth,c.StatMana,c.StatStamina,c.StatIntellect,c.StatSpirit,c.StatSpellPower,c.StatSpellCrit,c.StatSpellHaste,c.StatMP5],defaults:{gear:F.gear,epWeights:u.fromMap({[c.StatIntellect]:.38,[c.StatSpirit]:.34,[c.StatSpellPower]:1,[c.StatSpellCrit]:.69,[c.StatSpellHaste]:.77,[c.StatMP5]:0}),consumes:M,rotation:k,talents:B.data,specOptions:H,raidBuffs:f.create({giftOfTheWild:g.TristateEffectImproved,powerWordFortitude:g.TristateEffectImproved,strengthOfEarthTotem:g.TristateEffectImproved,arcaneBrilliance:!0,unleashedRage:!0,leaderOfThePack:g.TristateEffectRegular,icyTalons:!0,totemOfWrath:!0,demonicPact:500,swiftRetribution:!0,moonkinAura:g.TristateEffectRegular,sanctifiedRetribution:!0,manaSpringTotem:g.TristateEffectRegular,bloodlust:!0,thorns:g.TristateEffectImproved,devotionAura:g.TristateEffectImproved,shadowProtection:!0}),partyBuffs:h.create({}),individualBuffs:p.create({blessingOfKings:!0,blessingOfSanctuary:!0,blessingOfWisdom:g.TristateEffectImproved,blessingOfMight:g.TristateEffectImproved}),debuffs:S.create({judgementOfWisdom:!0,judgementOfLight:!0,misery:!0,faerieFire:g.TristateEffectImproved,ebonPlaguebringer:!0,totemOfWrath:!0,shadowMastery:!0,bloodFrenzy:!0,mangle:!0,exposeArmor:!0,sunderArmor:!0,vindication:!0,thunderClap:g.TristateEffectImproved,insectSwarm:!0})},playerIconInputs:[],rotationInputs:T,includeBuffDebuffInputs:[],excludeBuffDebuffInputs:[],otherInputs:{inputs:[y,I,E,w]},encounterPicker:{showExecuteProportion:!1},presets:{talents:[B],gear:[A,F,j]},autoRotation:e=>O.create()})}}export{k as D,G as H,F as P,B as S,H as a,M as b,j as c};
//# sourceMappingURL=sim-98225bb2.chunk.js.map
