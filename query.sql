SELECT * FROM health_cares 
    WHERE 
    heart_disease=true AND 
    bmi = (SELECT ROUND(avgBmi::numeric,0) FROM
     (SELECT AVG(bmi) as avgBmi FROM health_cares WHERE sex='Female' 
     AND age_category=(SELECT MIN(age_category) FROM health_cares) 
     AND sleep_time=(SELECT ROUND(avgSleep::numeric,0) FROM (SELECT AVG(sleep_time) as avgSleep FROM health_cares) as bmi_avg_t)) as avgBmiT) 
     AND 
    sleep_time = (SELECT ROUND(avg_sleep::numeric,0) FROM
     (SELECT AVG(sleep_time) as avg_sleep FROM health_cares WHERE sex='Female' AND 
     bmi=(SELECT AVG(avg_bmi) FROM (SELECT race, ROUND(avg_bmi::numeric,0) as avg_bmi FROM (SELECT race,AVG(bmi) as avg_bmi 
     FROM health_cares GROUP BY race) as complex_query_cond2) as round_AvgSleepT)) as roundAvgSleepT) 
     AND 
    sex= 
    (SELECT sex FROM(SELECT COUNT(sex) as numberPP,sex FROM health_cares GROUP BY sex) as numberPPT 
    WHERE numberPP= (SELECT MAX(totalPP) FROM (SELECT COUNT(sex) as totalPP,sex FROM health_cares GROUP BY sex) as maxPPT))
     AND  
    age_category = 
    (SELECT MAX(DISTINCT age_category) FROM health_cares WHERE bmi = (SELECT ROUND(avgBmi,0)
     FROM (SELECT AVG(bmi) as avgBmi FROM health_cares) as avgBmiT)) 
    AND 
    physical_health = (SELECT ROUND(avgPh::numeric,0) 
    FROM (SELECT AVG(physical_health) as avgPh FROM health_cares) as roundPH) ORDER BY index