import{fb as e,A as t,fc as a,fd as s,fe as i,ff as l,cg as o,ch as n,e1 as r,L as d,p as c,q as h,fg as p,E as S,cp as m,cq as u,cr as f,ct as g,cs as I,a2 as v,ab as C,F as T,aG as b,br as O,w as y,B as P,aH as A}from"./detailed_results-2fca7066.chunk.js";import{a as R,y as k,m as w,o as D,d as B,c as E,b as M,e as x,f as N,I as G,W as H,T as F,_ as j,$ as W,a0 as V,a1 as L,a3 as z,a2 as _,x as q}from"./preset_utils-c298d622.chunk.js";const U=R({fieldName:"startingRage",label:"Starting Rage",labelTooltip:"Initial rage at the start of each iteration."}),K=k({fieldName:"shout",values:[{color:"c79c6e",value:e.WarriorShoutNone},{actionId:t.fromSpellId(47436),value:e.WarriorShoutBattle},{actionId:t.fromSpellId(469),value:e.WarriorShoutCommanding}]}),$=w({fieldName:"useShatteringThrow",id:t.fromSpellId(64382)}),J={inputs:[D({fieldName:"customRotation",numColumns:3,values:[{actionId:t.fromSpellId(57823),value:a.Revenge},{actionId:t.fromSpellId(47488),value:a.ShieldSlam},{actionId:t.fromSpellId(47440),value:a.Shout},{actionId:t.fromSpellId(47502),value:a.ThunderClap},{actionId:t.fromSpellId(25203),value:a.DemoralizingShout},{actionId:t.fromSpellId(47486),value:a.MortalStrike},{actionId:t.fromSpellId(47498),value:a.Devastate},{actionId:t.fromSpellId(47467),value:a.SunderArmor},{actionId:t.fromSpellId(12809),value:a.ConcussionBlow},{actionId:t.fromSpellId(46968),value:a.Shockwave}]}),B({fieldName:"hsRageThreshold",label:"HS rage threshold",labelTooltip:"Heroic Strike when rage is above:"}),E({fieldName:"prioSslamOnShieldBlock",label:"Prio SSlam on Shield Block",labelTooltip:"The rotation code will prio SSlam over Revenge during active shield block windows."}),M({fieldName:"demoShoutChoice",label:"Demo Shout",values:[{name:"Never",value:s.DemoShoutChoiceNone},{name:"Maintain Debuff",value:s.DemoShoutChoiceMaintain},{name:"Filler",value:s.DemoShoutChoiceFiller}]}),M({fieldName:"thunderClapChoice",label:"Thunder Clap",values:[{name:"Never",value:i.ThunderClapChoiceNone},{name:"Maintain Debuff",value:i.ThunderClapChoiceMaintain},{name:"On CD",value:i.ThunderClapChoiceOnCD}]})]},Q={items:[{id:40546,enchant:3818,gems:[41380,40034]},{id:40387},{id:39704,enchant:3852,gems:[40034]},{id:40722,enchant:3605},{id:44e3,enchant:3832,gems:[40034,40015]},{id:39764,enchant:3850,gems:[0]},{id:40545,enchant:3860,gems:[40034,0]},{id:39759,enchant:3601,gems:[40008,36767]},{id:40589,enchant:3822},{id:39717,enchant:3232,gems:[40089]},{id:40370},{id:40718},{id:40257},{id:44063,gems:[36767,40089]},{id:40402,enchant:3788},{id:40400,enchant:3849},{id:41168,gems:[36767]}]},X={items:[{id:46166,enchant:3818,gems:[41380,40008]},{id:45485,gems:[40008]},{id:46167,enchant:3852,gems:[40008]},{id:45496,enchant:3605,gems:[40023]},{id:46162,enchant:3832,gems:[40008,40008]},{id:45111,enchant:3850,gems:[0]},{id:45487,enchant:3860,gems:[40008,40008,0]},{id:45139,enchant:3601,gems:[40008]},{id:46169,enchant:3822,gems:[40088,40008]},{id:45988,enchant:3232,gems:[36767,36767]},{id:45471,gems:[45880]},{id:45247},{id:45158},{id:46021},{id:45442,enchant:3788,gems:[40034]},{id:45587,enchant:3849,gems:[36767]},{id:45137,enchant:3608}]},Y={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:47440}}},doAtValue:{const:{val:"-10s"}}},{action:{castSpell:{spellId:{otherId:"OtherActionPotion"}}},doAtValue:{const:{val:"-1s"}}}],priorityList:[{action:{schedule:{schedule:"29s, 209s",innerAction:{castSpell:{spellId:{spellId:12975}}}}}},{action:{condition:{cmp:{op:"OpGe",lhs:{currentRage:{}},rhs:{const:{val:"30"}}}},castSpell:{spellId:{tag:1,spellId:47450}}}},{action:{autocastOtherCooldowns:{}}},{action:{castSpell:{spellId:{spellId:47488}}}},{action:{castSpell:{spellId:{spellId:57823}}}},{action:{condition:{auraShouldRefresh:{sourceUnit:{type:"Self"},auraId:{spellId:47440},maxOverlap:{const:{val:"3s"}}}},castSpell:{spellId:{spellId:47440}}}},{action:{condition:{auraShouldRefresh:{auraId:{spellId:47502},maxOverlap:{const:{val:"2s"}}}},castSpell:{spellId:{spellId:47502}}}},{action:{condition:{auraShouldRefresh:{auraId:{spellId:47437},maxOverlap:{const:{val:"2s"}}}},castSpell:{spellId:{spellId:25203}}}},{action:{castSpell:{spellId:{spellId:47498}}}}]},Z=x("PreRaid Balanced",{items:[{id:42549,enchant:3818,gems:[41380,40015]},{id:40679},{id:37814,enchant:3852},{id:37728,enchant:3605},{id:39611,enchant:1953,gems:[40008,40008]},{id:37620,enchant:3850,gems:[0]},{id:39622,enchant:3860,gems:[40034,0]},{id:37379,enchant:3601,gems:[40034,36767]},{id:43500,enchant:3822,gems:[40034]},{id:44201,enchant:3232},{id:37784},{id:37186},{id:37220},{id:44063,gems:[36767,40089]},{id:37401,enchant:3788},{id:43085,enchant:3849},{id:41168,gems:[36767]}]}),ee=x("P1 Balanced Preset",Q),te=x("P2 Survival Preset",X),ae=l.create({customRotation:o.create({spells:[n.create({spell:a.ShieldSlam}),n.create({spell:a.Revenge}),n.create({spell:a.Shout}),n.create({spell:a.ThunderClap}),n.create({spell:a.DemoralizingShout}),n.create({spell:a.MortalStrike}),n.create({spell:a.Devastate}),n.create({spell:a.SunderArmor}),n.create({spell:a.ConcussionBlow}),n.create({spell:a.Shockwave})]}),demoShoutChoice:s.DemoShoutChoiceNone,thunderClapChoice:i.ThunderClapChoiceNone,hsRageThreshold:30}),se=N("Default",Y),ie={name:"Standard",data:r.create({talentsString:"2500030023-302-053351225000012521030113321",glyphs:d.create({major1:c.GlyphOfBlocking,major2:c.GlyphOfVigilance,major3:c.GlyphOfDevastate,minor1:h.GlyphOfCharge,minor2:h.GlyphOfThunderClap,minor3:h.GlyphOfCommand})})},le={name:"UA",data:r.create({talentsString:"35023301230051002020120002-2-05035122500000252",glyphs:d.create({major1:c.GlyphOfRevenge,major2:c.GlyphOfHeroicStrike,major3:c.GlyphOfSweepingStrikes,minor1:h.GlyphOfCharge,minor2:h.GlyphOfThunderClap,minor3:h.GlyphOfCommand})})},oe=p.create({shout:e.WarriorShoutCommanding,useShatteringThrow:!1,startingRage:0}),ne=S.create({battleElixir:m.ElixirOfExpertise,guardianElixir:u.ElixirOfProtection,food:f.FoodDragonfinFilet,defaultPotion:g.IndestructiblePotion,prepopPotion:g.IndestructiblePotion,thermalSapper:!0,fillerExplosive:I.ExplosiveSaroniteBomb});class re extends G{constructor(e,t){super(e,t,{cssClass:"protection-warrior-sim-ui",cssScheme:"warrior",knownIssues:[],epStats:[v.StatStamina,v.StatStrength,v.StatAgility,v.StatAttackPower,v.StatExpertise,v.StatMeleeHit,v.StatMeleeCrit,v.StatMeleeHaste,v.StatArmor,v.StatBonusArmor,v.StatArmorPenetration,v.StatDefense,v.StatBlock,v.StatBlockValue,v.StatDodge,v.StatParry,v.StatResilience,v.StatNatureResistance,v.StatShadowResistance,v.StatFrostResistance],epPseudoStats:[C.PseudoStatMainHandDps],epReferenceStat:v.StatAttackPower,displayStats:[v.StatHealth,v.StatArmor,v.StatBonusArmor,v.StatStamina,v.StatStrength,v.StatAgility,v.StatAttackPower,v.StatExpertise,v.StatMeleeHit,v.StatMeleeCrit,v.StatMeleeHaste,v.StatArmorPenetration,v.StatDefense,v.StatBlock,v.StatBlockValue,v.StatDodge,v.StatParry,v.StatResilience,v.StatNatureResistance,v.StatShadowResistance,v.StatFrostResistance],defaults:{gear:te.gear,epWeights:T.fromMap({[v.StatArmor]:.174,[v.StatBonusArmor]:.155,[v.StatStamina]:2.336,[v.StatStrength]:1.555,[v.StatAgility]:2.771,[v.StatAttackPower]:.32,[v.StatExpertise]:1.44,[v.StatMeleeHit]:1.432,[v.StatMeleeCrit]:.925,[v.StatMeleeHaste]:.431,[v.StatArmorPenetration]:1.055,[v.StatBlock]:1.32,[v.StatBlockValue]:1.373,[v.StatDodge]:2.606,[v.StatParry]:2.649,[v.StatDefense]:3.305},{[C.PseudoStatMainHandDps]:6.081}),consumes:ne,rotation:ae,talents:ie.data,specOptions:oe,raidBuffs:b.create({giftOfTheWild:O.TristateEffectImproved,powerWordFortitude:O.TristateEffectImproved,abominationsMight:!0,swiftRetribution:!0,bloodlust:!0,strengthOfEarthTotem:O.TristateEffectImproved,leaderOfThePack:O.TristateEffectImproved,sanctifiedRetribution:!0,devotionAura:O.TristateEffectImproved,stoneskinTotem:O.TristateEffectImproved,icyTalons:!0,retributionAura:!0,thorns:O.TristateEffectImproved,shadowProtection:!0}),partyBuffs:y.create({}),individualBuffs:P.create({blessingOfKings:!0,blessingOfMight:O.TristateEffectImproved,blessingOfSanctuary:!0}),debuffs:A.create({sunderArmor:!0,mangle:!0,vindication:!0,faerieFire:O.TristateEffectImproved,insectSwarm:!0,bloodFrenzy:!0,judgementOfLight:!0,heartOfTheCrusader:!0,frostFever:O.TristateEffectImproved})},playerIconInputs:[K,$],rotationInputs:J,includeBuffDebuffInputs:[H],excludeBuffDebuffInputs:[],otherInputs:{inputs:[F,j,W,V,L,z,_,U,q]},encounterPicker:{showExecuteProportion:!1},presets:{talents:[ie,le],rotations:[se],gear:[Z,ee,te]},autoRotation:e=>se.rotation.rotation})}}export{ae as D,re as P,ie as S,oe as a,ne as b,ee as c,te as d};
//# sourceMappingURL=sim-2415f8b2.chunk.js.map
